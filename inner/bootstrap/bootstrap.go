package innerbootstrap

import (
	"fmt"
	"os"
	"os/exec"
	"ekathvam/shared"
)

func Bootstrap() {
	fmt.Println("Identifying running environment")
	distro, err := globalshared.OS_identify("/etc/os-release")
	if err != nil {
		fmt.Println("OS identification failed:", err)
		os.Exit(1)
	}
	if distro == "debian" {
		fmt.Println("\nDebian detected successfully.")
	} else {
		fmt.Println("\nUnsupported Container or running outside the container:", distro)
		os.Exit(1)
	}

	fmt.Println("Verifying Superuser Privileage...")
	if os.Geteuid() == 0 {
		fmt.Println("Running with root privileges!")
	} else {
		fmt.Println("OOPS! Bootstrap running as a regular user.\nExiting....")
		os.Exit(1)
	}
	
	fmt.Println("Installing necessary packages inside the container!...")
	globalshared.Install_pkg("ufw")

	update := exec.Command("apt", "update", "&&", "apt", "upgrade", "-y")
	globalshared.Cmd_run("Updating the container environment!", update)

	ufw := exec.Command("ufw", "enable")
	globalshared.Cmd_run("Enabling the UFW", ufw)




	var option string
	fmt.Println("\nWelcome to Ekathvam Inner Bootstrap!")
	fmt.Println("\nWhat would you like this machine to provide?\n\n1. Administrator\n2. Client\n3. Web Server\n4. DNS Server\n\nSelect an option:")
	for {
		fmt.Print("Enter option (1-4): ")
		fmt.Scanln(&option)

		switch option {
		case "1":
			fmt.Println("Downloading Administrator Server Bootstrap....")
			err = globalshared.Install_pkg("https://raw.githubusercontent.com/ekathvam/ekathvam/main/bin/administrator_bootstrap")
			if err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}

			cmd := exec.Command("/ekathvam/downloads/administrator_bootstrap")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
				
			err = cmd.Run()
			
			if err != nil {
			    fmt.Fprintf(os.Stderr, "Administrator bootstrap failed to launch: %v\n", err)
			    os.Exit(1)
			}
			
			fmt.Println("\n[Host] Inner bootstrap completed.")

		case "2":
			fmt.Println("Downloading Client Server Bootstrap....")
			err = globalshared.Install_pkg("https://raw.githubusercontent.com/ekathvam/ekathvam/main/bin/client_bootstrap")
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			cmd := exec.Command("/ekathvam/downloads/client_bootstrap")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
				
			err = cmd.Run()
			
			if err != nil {
			    fmt.Fprintf(os.Stderr, "Client bootstrap failed to launch: %v\n", err)
			    os.Exit(1)
			}
			
			fmt.Println("\n[Host] Inner bootstrap completed.")

		case "3":
			fmt.Println("Downloading Web Server Bootstrap....")
			err = globalshared.Install_pkg("https://raw.githubusercontent.com/ekathvam/ekathvam/main/bin/web_bootstrap")
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			cmd := exec.Command("/ekathvam/downloads/web_bootstrap")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
				
			err = cmd.Run()
			
			if err != nil {
			    fmt.Fprintf(os.Stderr, "Web bootstrap failed to launch: %v\n", err)
			    os.Exit(1)
			}
			
			fmt.Println("\n[Host] Inner bootstrap completed.")

		case "4":
			fmt.Println("Downloading DNS Bootstrap....")
			err = globalshared.Install_pkg("https://raw.githubusercontent.com/ekathvam/ekathvam/main/bin/dns_bootstrap")
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			cmd := exec.Command("/ekathvam/downloads/dns_bootstrap")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
				
			err = cmd.Run()
			
			if err != nil {
			    fmt.Fprintf(os.Stderr, "DNS bootstrap failed to launch: %v\n", err)
			    os.Exit(1)
			}
			
			fmt.Println("\n[Host] Inner bootstrap completed.")

		default:
			fmt.Println("❌ Invalid selection. Please choose a valid option.")
		}
	}
}



