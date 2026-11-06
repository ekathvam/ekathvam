package globalshared

import (
	"fmt"
	"os"
	"os/exec"
	"net/url"
	"net/http"
	"io"
	"path"
	"path/filepath"
	"strings"
	"errors"
)

/*	Download Types
	1. Linux Packages
	2. URL
*/

func Install_pkg(pkg string) error {
	
	pkg_type := pkgtypeparser(pkg)
	if pkg_type == "1" {
		err := linux_pkg(pkg)
		return err
	} else {
		err := url_pkg(pkg)
		return err
	}

}

func linux_pkg(pkg_name string) error {

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

 func url_pkg(fileURL string) error {
	// 1. Parse the URL string to cleanly extract its path component
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return fmt.Errorf("invalid URL structure: %w", err)
	}

	// 2. Extract the absolute last part of the path string
	fileName := path.Base(parsedURL.Path)
	if fileName == "" || fileName == "." || fileName == "/" {
		return errors.New("the URL string does not terminate with a valid file name segment")
	}

	// 3. Define the mandatory target directory and join with the filename
	targetDir := "/ekathvam/downloads"
	finalSavePath := filepath.Join(targetDir, fileName)

	// 4. Ensure the parent directory structure exists with correct permissions
	err = os.MkdirAll(targetDir, 0700)
	if err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// 5. Send the HTTP GET request
	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad server status: %s", resp.Status)
	}

	// 6. Create the local file at the absolute target path
	out, err := os.Create(finalSavePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 7. Stream the download straight to disk
	fmt.Printf("Downloading and auto-saving to: %s\n", finalSavePath)
	_, err = io.Copy(out, resp.Body)
	return err
 }


func pkgtypeparser(input string) string {
	// Clean up whitespace
	input = strings.TrimSpace(input)

	// Attempt to parse the string as a URL
	u, err := url.Parse(input)
	
	// If it fails basic parsing, it's definitely just a string
	if err != nil {
		return "1"
	}

	// Crucial Check: A valid web URL MUST have a scheme (http/https) and a host.
	// Without this, plain text like "my-file.txt" or "hello" gets misidentified as a URL path.
	if (u.Scheme == "http" || u.Scheme == "https") && u.Host != "" {
		return "2"
	}

	return "1"
}
