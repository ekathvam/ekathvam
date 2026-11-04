package main

import (
	"fmt"
	"os"
	"ekathvam/shared"
	"ekathvam/core/dns"
	"ekathvam/core/nginx"
)

func main() {

	var option int

	fmt.Println("Hello welcome to Ekathvam!")
	if os.Geteuid() != 0 {
		fmt.Println("Run with root privileges (sudo)\nExiting....")
		os.Exit(1)
	}
	err := shared.Create_dir("/ekathvam")
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Print("\n1.Configure and host a website.\n2.Configure and host a DNS server.\n Option : ")
	fmt.Scanln(&option)

	if option == 1 {
		nginx.Nginx()
	} else if option == 2 {
		dns.Dns()
	} else {
		fmt.Println("Invaild Option!")
	}
}

