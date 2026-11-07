package dnsbootstrap

import (
	"fmt"
	"ekathvam/shared"
	"os"
	"os/exec"
	"net"
	"time"
	"bufio"
)

func Dnsbootstrap() {
	fmt.Println("Initializing DNS Server Bootstrap..")
	fmt.Println("Installing NGINX")
	fmt.Println("Verifying Superuser Privileage...")
	if os.Geteuid() == 0 {
		fmt.Println("Running with root privileges!")
	} else {
		fmt.Println("OOPS! Bootstrap running as a regular user.\nExiting....")
		os.Exit(1)
	}
	
	pkg := []string{"pdns-server", "pdns-backend-sqlite3", "sqlite3"}
	for _, pkg := range pkg {
		err := globalshared.Install_pkg(pkg)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
	}


	
	cmd3 := exec.Command("ufw", "delete", "allow", "53/tcp", "ufw", "delete", "allow", "2008/tcp",
		"&&", "ufw", "allow", "from", "192.168.1.0/24", "to", "any", "port", "53", "proto", "tcp",
		"&&", "ufw", "allow", "from", "192.168.1.0/24", "to", "any", "port", "2008", "proto", "tcp",
		"&&", "ufw", "reload",
	)
	globalshared.Cmd_run("Unlocking firwalls", cmd3)
	
	cmd1 := exec.Command("systemctl", "enable","pdns")
	globalshared.Cmd_run("Enabling Nginx", cmd1)

	cmd2 := exec.Command("systemctl", "start", "pdns")
	globalshared.Cmd_run("Starting Nginx", cmd2)


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

}


