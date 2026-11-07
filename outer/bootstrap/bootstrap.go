package outerbootstrap

import (
	"fmt"
	"os"
	"os/exec"
	"ekathvam/shared"
	"strings"
	"time"
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
	cmd3 := exec.Command("lxc-start", "-n", "ekathvam", "-P", "/ekathvam/containers")
	globalshared.Cmd_run("Starting LXC Container", cmd3)

	containerName := "ekathvam"
	storagePath := "/ekathvam/containers"

	
	containerRunning := false

	for i := 0; i < 10; i++ {
    	cmd := exec.Command(
        	"lxc-info",
	        "-n", containerName,
	        "-P", storagePath,
	        "-s",
	    )
	
	    output, err := cmd.Output()
	
	    if err == nil && strings.Contains(string(output), "RUNNING") {
	        containerRunning = true
	        break
	    }
	
	    time.Sleep(500 * time.Millisecond)
	}

	if !containerRunning {
	    fmt.Println("Container failed to reach RUNNING state.")
	    os.Exit(1)
	}



	fmt.Println("\n Container Has been started")
	
	fmt.Println("Downloading Inner Core Bootstrap.....")
	err = globalshared.Install_pkg("https://raw.githubusercontent.com/ekathvam/ekathvam/main/bin/inner_bootstrap")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	
	fmt.Println("Creating downloads directory...")
	err = os.MkdirAll("/ekathvam/downloads", 0700)
	if err != nil {
		  fmt.Printf("failed to create directory structure: %w", err)
	}
		err = os.MkdirAll("/ekathvam/containers/ekathvam/rootfs/ekathvam/downloads/", 0700)
	if err != nil {
		  fmt.Printf("failed to create directory structure: %w", err)
	}

	
	mv_inner_bootstrap := exec.Command("mv", "/ekathvam/downloads/inner_bootstrap", "/ekathvam/containers/ekathvam/rootfs/ekathvam/downloads/")
	globalshared.Cmd_run("Moving Inner Core Bootstrap to the Container..", mv_inner_bootstrap)


	fmt.Println("⚙️  [Entering Ekathvam Container]")
	
	cmd := exec.Command(
	    "lxc-attach",
	    "-n", containerName,
	    "-P", storagePath,
	    "--",
	    "/bin/bash",
	    "-c",
	    "chmod +x /ekathvam/downloads/inner_bootstrap && exec /ekathvam/downloads/inner_bootstrap",
	)
	
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	err = cmd.Run()
	
	if err != nil {
	    fmt.Fprintf(os.Stderr, "Inner bootstrap failed to launch: %v\n", err)
	    os.Exit(1)
	}
	
	fmt.Println("\n[Host] Outer bootstrap completed.")

}


