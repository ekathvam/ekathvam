package globalshared

import (
	"fmt"
	"os/exec"
)


func Install_pkg(pkg_name string) error {
	var cmd *exec.Cmd

	pm := get_pm()
	if pm == "unknown" {
		return fmt.Errorf("failed to identify the package manager.")
	}

	switch pm {
	case "apt":
		cmd = exec.Command("apt-get", "install", "-y", pkg_name)
	case "dnf":
		cmd = exec.Command("dnf", "install", "-y", pkg_name)
	case "pacman":
		cmd = exec.Command("pacman", "-S", "--noconfirm", "--needed", pkg_name)
	case "zypper":
		cmd = exec.Command("zypper", "--non-interactive", "install", pkg_name)
	case "apk":
		cmd = exec.Command("apk", "add", pkg_name)
	case "pkg": // FreeBSD style package manager used on some Linux/Termux systems
		cmd = exec.Command("pkg", "install", "-y", pkg_name)
	case "yum":
		cmd = exec.Command("yum", "install", "-y", pkg_name)
		
	// --- Bonus Package Managers You Missed ---
	case "brew": // Homebrew for Linux (increasingly common on Ubuntu/Fedora)
		cmd = exec.Command("brew", "install", pkg_name)
	case "swupd": // Intel's Clear Linux OS bundle manager
		cmd = exec.Command("swupd", "bundle-add", pkg_name)
	case "nix-env": // Nix package manager (used on NixOS and standalone)
		cmd = exec.Command("nix-env", "-iA", pkg_name)
		
	default:
		return fmt.Errorf("unsupported package manager: %s", pm)
	}

	// Execute the command and capture output/errors
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install %s using %s: %v\nOutput: %s", pkg_name, pm, err, string(output))
		
	}

	fmt.Printf("Successfully processed %s via %s\n", pkg_name, pm)
	return nil
}

func get_pm () string {
	fmt.Println("Analysing the pacakage manager.")
	for _, pm := range []string{"apt","dnf","pacman","zypper","apk","pkg","yum"} {
		if _, err := exec.LookPath(pm); err == nil {
			return pm
		}
	}
	return "unknown"
}


