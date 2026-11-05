package globalshared

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//Templates

/*
TEXT_APPEND_ROOT
TEXT_APPEND_READ
TEXT_OVERWRITE_ROOT
TEXT_OVERWRITE_READ

FILE_APPEND_ROOT
FILE_APPEND_READ
FILE_OVERWRITE_ROOT
FILE_OVERWRITE_READ
*/

// Content_writer routes and processes data based on a 3-part template combination:
// Format (TEXT/FILE) + Action (APPEND/OVERWRITE) + Permission (ROOT/READ)
func Content_writer(source, destin, template string) error {
	// 1. Parse template configurations
	flag, perm, err := parseTemplate(template)
	if err != nil {
		return err
	}

	// 2. Prepare the destination file safely (creates directories and opens file)
	destFile, err := prepareDestinFile(destin, flag, perm)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// 3. Perform data operation based on format type
	if strings.HasPrefix(template, "TEXT_") {
		_, err = destFile.WriteString(source)
	} else if strings.HasPrefix(template, "FILE_") {
		err = streamCopy(source, destFile)
	} else {
		err = fmt.Errorf("unknown template format: %s", template)
	}

	if err != nil {
		return fmt.Errorf("failed to write content to destination: %w", err)
	}
	return nil
}

// prepareDestinFile isolates directory creation and file opening so it is written exactly once
func prepareDestinFile(destin string, flag int, perm os.FileMode) (*os.File, error) {
	destin_dir := filepath.Dir(destin)
	if err := os.MkdirAll(destin_dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create destination directories: %w", err)
	}

	file, err := os.OpenFile(destin, flag, perm)
	if err != nil {
		return nil, fmt.Errorf("failed to open/create destination file: %w", err)
	}
	return file, nil
}

// streamCopy handles opening the source path and streaming its payload to the destination
func streamCopy(sourcePath string, destFile *os.File) error {
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

// parseTemplate extracts systemic flags and permissions from the template string
func parseTemplate(template string) (int, os.FileMode, error) {
	var flag int
	var perm os.FileMode

	// Determine Action
	if strings.Contains(template, "_APPEND_") {
		flag = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	} else if strings.Contains(template, "_OVERWRITE_") {
		flag = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	} else {
		return 0, 0, fmt.Errorf("unknown template action: %s", template)
	}

	// Determine Permissions
	if strings.HasSuffix(template, "_ROOT") {
		perm = 0700
	} else if strings.HasSuffix(template, "_READ") {
		perm = 0644
	} else {
		return 0, 0, fmt.Errorf("unknown template permission: %s", template)
	}

	return flag, perm, nil
}

