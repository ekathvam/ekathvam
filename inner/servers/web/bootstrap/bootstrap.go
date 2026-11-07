package webbootstrap

import (
	"fmt"
	"ekathvam/shared"
	"os"
	"os/exec"
	"net"
	"encoding/json"
	"path/filepath"
	"bufio"
	"log"
	"regexp"
	"time"
)

func Webbootstrap() {
	fmt.Println("Initializing Web Server Bootstrap..")
	fmt.Println("Verifying Superuser Privileage...")
	if os.Geteuid() == 0 {
		fmt.Println("Running with root privileges!")
	} else {
		fmt.Println("OOPS! Bootstrap running as a regular user.\nExiting....")
		os.Exit(1)
	}

	err := globalshared.Install_pkg("nginx")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	
	cmd3 := exec.Command("ufw", "delete", "allow", "80/tcp",
		"&&", "ufw", "delete", "allow", "443/tcp",
		"&&", "ufw", "allow", "from", "192.168.1.0/24", "to", "any", "port", "80", "proto", "tcp",
		"&&", "ufw", "allow", "from", "192.168.1.0/24", "to", "any", "port", "443", "proto", "tcp",
		"&&", "ufw", "allow", "from", "192.168.1.0/24", "to", "any", "port", "2008", "proto", "tcp",
		"&&", "ufw", "reload",
	)
	globalshared.Cmd_run("Unlocking firwalls", cmd3)
	
	cmd1 := exec.Command("systemctl", "enable","nginx")
	globalshared.Cmd_run("Enabling Nginx", cmd1)

	cmd2 := exec.Command("systemctl", "start", "nginx")
	globalshared.Cmd_run("Starting Nginx", cmd2)
	
	var domain string
	for {	

		fmt.Printf("What is the domain that you wanted to resolved to this page? (like example.com,google.com,github.org: )")
		fmt.Scanln(&domain)
		
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`, domain)
	
		if !matched {
				fmt.Printf("Invalid format: %s   (enter vaild domain like ekathvam.eka or example.com)", domain)
		} else {
			break
		}
	}

	// 8.8.8.8:80 does not need to be reachable, this is a dummy connection
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	// 1. Your data (can be a map, struct, etc.)
	data := map[string]interface{}{
		"DOMAIN": domain,
		"IP":  localAddr.IP,
	}

	// 2. Specify your path and file name
	
	filePath := "/path/to/your/folder/output.json"


	dir := filepath.Dir(filePath)

	// Create the directories if they don't exist (0755 grants standard folder permissions)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Println("Error creating directories:", err)
		return

	}

	// 3. Format the JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}

	// 4. Save to the specific location (0644 gives read/write permissions)
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}

	fmt.Println("JSON successfully saved to:", filePath)

for {
		// 1. Find the local LAN or loopback IP address
		var ip net.IP
		addrs, _ := net.InterfaceAddrs()
		for _, a := range addrs {
			if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsUnspecified() && ipNet.IP.To4() != nil {
				if ipNet.IP.IsPrivate() || ipNet.IP.IsLoopback() {
					ip = ipNet.IP
					if ipNet.IP.IsPrivate() {
						break
					}
				}
			}
		}
		if ip == nil {
			ip = net.IPv4(127, 0, 0, 1)
		}

		target := fmt.Sprintf("%s:2008", ip.String())

		// 2. Dial the target server
		conn, err := net.DialTimeout("tcp", target, 2*time.Second)
		if err != nil {
			fmt.Printf("[%s] Reconnecting to %s...\n", time.Now().Format("03:04:05 PM"), target)
			time.Sleep(3 * time.Second)
			continue
		}

		// 3. Send the outbound question
		_, err = conn.Write([]byte("Admin are you there?\n"))
		if err != nil {
			conn.Close()
			time.Sleep(3 * time.Second)
			continue
		}
		fmt.Printf("[%s] Sent query to %s\n", time.Now().Format("03:04:05 PM"), target)

		// 4. Set a read deadline and listen for the confirmation response
		conn.SetReadDeadline(time.Now().Add(4 * time.Second))
		reply, err := bufio.NewReader(conn).ReadString('\n')
		
		if err == nil && reply == "Admin is here!\n" {
			fmt.Printf("[%s] Received authorization confirmation: %sClosing link cleanly.\n", time.Now().Format("03:04:05 PM"), reply)
			conn.Close()
			return // Exits the entire function and script execution stops immediately
		}

		// Clean drop if no reply or incorrect response matches
		conn.Close()
		time.Sleep(3 * time.Second)
	}

	fmt.Print("Your pages are live! place your necessary files into /var/nginx/www/html")
}


