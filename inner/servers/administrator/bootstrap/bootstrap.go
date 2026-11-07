package adminbootstrap

import (
	"fmt"
	"ekathvam/shared"
	"os"
	"os/exec"
	"net"
	"time"
	"bufio"
)

func Adminbootstrap() {
	fmt.Println("Initializing  Server Bootstrap..")
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


	
	cmd3 := exec.Command("ufw", "delete", "allow", "2008/tcp",
		"&&", "ufw", "allow", "from", "192.168.1.0/24", "to", "any", "port", "2008", "proto", "tcp",
		"&&", "ufw", "reload",
	)
	globalshared.Cmd_run("Unlocking firwalls", cmd3)
	


	// 1. Resolve local LAN interface or fallback to loopback
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

	// 2. Bind and listen on port 2008
	listener, err := net.Listen("tcp", target)
	if err != nil {
		fmt.Printf("Error starting admin listener: %v\n", err)
		return
	}
	defer listener.Close()
	fmt.Printf("[%s] Admin server listening on %s...\n", time.Now().Format("03:04:05 PM"), target)

	// 3. Execution loop to wait for incoming connections
	for {
		// This blocks until a client connects, safely instantiating 'conn' each iteration
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		// Enforce a strict 4-second read deadline for the packet arrival
		conn.SetReadDeadline(time.Now().Add(4 * time.Second))
		
		// Read the inbound text payload up to the newline character
		announcement, err := bufio.NewReader(conn).ReadString('\n')

		// Challenge validation logic
		expectedMsg := "Admiin are you there?\n"
		if err == nil && announcement == expectedMsg {
			fmt.Printf("[%s] Announcement verified successfully.\n", time.Now().Format("03:04:05 PM"))
			
			// Send response back directly on the established channel
			_, err = conn.Write([]byte("Admin is here!\n"))
			if err != nil {
				fmt.Printf("Failed to transmit signature response: %v\n", err)
			}
			
			conn.Close()
			fmt.Println("Handshake routine executed. Stopping execution.")
			return // Gracefully terminates the entire application block
		}

		// Silent Drop: The signature didn't match, or a timeout occurred.
		// Close the socket immediately without sending any network response back.
		conn.Close()
	}
}

