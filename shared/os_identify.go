package globalshared

import (
	"fmt"
	"bufio"
	"strings"
	"os"
	
)

type OSInfo struct {
	Name       string
	ID         string
	VersionID  string
	Version    string
	Codename   string
	PrettyName string
}

func OS_identify(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "error", fmt.Errorf("failed to open %s: %w", path, err)
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
		return "error", fmt.Errorf("failed to read %s: %w", path, err)
	}

	info := OSInfo{
		Name:       values["NAME"],
		ID:         values["ID"],
		VersionID:  values["VERSION_ID"],
		Version:    values["VERSION"],
		Codename:   values["VERSION_CODENAME"],
		PrettyName: values["PRETTY_NAME"],
	}
	fmt.Println("Container System Details")
	fmt.Println("------------------------")
	fmt.Println("Name:", info.Name)
	fmt.Println("ID:", info.ID)
	fmt.Println("Version ID:", info.VersionID)
	fmt.Println("Version:", info.Version)
	fmt.Println("Codename:", info.Codename)
	fmt.Println("Pretty Name:", info.PrettyName)


	return info.ID, nil
}


