package outerbootstrap

import (
	"fmt"
	"os"
	"os/exec"
	"ekathvam/shared"
	)

func Bootstrap() {

fmt.Println("Verifying Superuser Privileage...")
	if os.Geteuid() == 0 {
		fmt.Println("Successfull. Running with root privileges!")
	} else {
		fmt.Println("Running as a regular user.\nRun with sudo or as root user.\nExiting....")
	}
	err := globalshared.Install_pkg("lxc")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	fmt.Println("🚀 Initializing LXC Deployment...")


	// 1. Create the LXC Container
	cmd1 := exec.Command("lxc-create",
		"-n", "ekathvam",
		"-P", "/ekathvam/containers",
		"-t", "download",
		"--",
		"--dist", "debian",
		"--release", "trixie",
		"--arch", "amd64",
	)
	globalshared.Cmd_run("Creating Container", cmd1)

	// 2. Enable lxc-net Service
	cmd2 := exec.Command("systemctl", "enable", "--now", "lxc-net")
	globalshared.Cmd_run("Enabling lxc-net service", cmd2)

	// 3. Start the Container
	cmd3 := exec.Command("lxc-start", "-n", "ekathvam", "-P", "/ekathvam/containers -d")
	globalshared.Cmd_run("Starting LXC Container", cmd3)

	fmt.Println("\n🎉 All steps completed successfully!")
}




