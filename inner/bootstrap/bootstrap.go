package innerbootstrap

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type OSInfo struct {
	Name       string
	ID         string
	VersionID  string
	Version    string
	Codename   string
	PrettyName string
}

func Bootstrap() {
	fmt.Println("Identifying running environment")
	info, err := ParseOSRelease("/etc/os-release")
	if err != nil {
		fmt.Println("OS identification failed:", err)
		os.Exit(1)
	}

	fmt.Println("Container System Details")
	fmt.Println("------------------------")
	fmt.Println("Name:", info.Name)
	fmt.Println("ID:", info.ID)
	fmt.Println("Version ID:", info.VersionID)
	fmt.Println("Version:", info.Version)
	fmt.Println("Codename:", info.Codename)
	fmt.Println("Pretty Name:", info.PrettyName)

	if info.ID == "debian" {
		fmt.Println("\nDebian detected successfully.")
	} else {
		fmt.Println("\nUnsupported Container or running outside the container:", info.ID)
		os.Exit(1)
	}

	fmt.Println("Verifying Superuser Privileage...")
	if os.Geteuid() == 0 {
		fmt.Println("Running with root privileges!")
	} else {
		fmt.Println("OOPS! Bootstrap running as a regular user.\nExiting....")
		os.Exit(1)
	}

	var option string
	fmt.Println("\nWelcome to Ekathvam Inner Bootstrap!")
	fmt.Println("\nWhat would you like this machine to provide?\n\n1. Administrator\n2. Client\n3. Web Server\n4. DNS Server\n\nSelect an option:")
	for {
		fmt.Print("Enter option (1-4): ")
		fmt.Scanln(&option)

		switch option {
		case "1":
			fmt.Println("Initiating Administrator Bootstrap...")
			os.Exit(0)
		case "2":
			fmt.Println("Initiating Client Bootstrap...")
			os.Exit(0)
		case "3":
			fmt.Println("Initiating Web Server Bootstrap...")
			os.Exit(0)
		case "4":
			fmt.Println("Initiating DNS Bootstrap...")
			os.Exit(0)
		default:
			fmt.Println("❌ Invalid selection. Please choose a valid option.")
		}
	}
}


func ParseOSRelease(path string) (OSInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return OSInfo{}, fmt.Errorf("failed to open %s: %w", path, err)
	}
	defer file.Close()

	values := make(map[string]string)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore blank lines and comments.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		// Remove surrounding quotes.
		value = strings.Trim(value, `"'`)

		values[key] = value
	}

	if err := scanner.Err(); err != nil {
		return OSInfo{}, fmt.Errorf("failed to read %s: %w", path, err)
	}

	info := OSInfo{
		Name:       values["NAME"],
		ID:         values["ID"],
		VersionID:  values["VERSION_ID"],
		Version:    values["VERSION"],
		Codename:   values["VERSION_CODENAME"],
		PrettyName: values["PRETTY_NAME"],
	}

	return info, nil
}


