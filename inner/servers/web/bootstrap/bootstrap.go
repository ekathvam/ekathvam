package webbootstrap

import (
	"fmt"
	"ekathvam/shared"
	"os"
	"os/exec"
)

func Webbootstrap() {
	fmt.Println("Initializing Web Server Bootstrap..")
	fmt.Println("Installing NGINX")
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
	cmd1 := exec.Command("systemctl", "enable","nginx")
	globalshared.Cmd_run("Enabling Nginx", cmd1)

	cmd2 := exec.Command("systemctl", "start", "nginx")
	globalshared.Cmd_run("Starting Nginx", cmd2)

	cmd3 := exec.Command("ufw", "allow", "&&", "ufw", "reload")
	globalshared.Cmd_run("Unlocking firwalls", cmd3)

	
}


