package globalshared

import (
	"fmt"
	"log"
	"os/exec"
	"bufio"
	"strings"
	"time"


)

func Cmd_run(stepName string, cmd *exec.Cmd) {
	// Intercept the standard output pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("[%s] Failed to create stdout pipe: %v", stepName, err)
	}

	// Direct standard error to standard output so everything streams together
	cmd.Stderr = cmd.Stdout

	// Start the command in the background
	if err := cmd.Start(); err != nil {
		log.Fatalf("[%s] Failed to start command: %v", stepName, err)
	}

	// Create a scanner to read the terminal lines in real-time
	scanner := bufio.NewScanner(stdout)
	
	fmt.Printf("\n⚙️  [%s]\n", stepName)
	
	// Create an animated indeterminate progress bar track
	// Uses \r (carriage return) to constantly overwrite the same line
	go func() {
		frames := []string{"[=   ]", "[==  ]", "[=== ]", "[ ===]", "[  ==]", "[   =]", "[    ]"}
		i := 0
		for cmd.ProcessState == nil {
			fmt.Printf("\r  %s Working... Progress Active ", frames[i%len(frames)])
			i++
			time.Sleep(150 * time.Millisecond)
		}
	}()

	// Read lines as they stream from the active Linux process
	for scanner.Scan() {
		line := scanner.Text()
		// Clean out extra empty space strings if necessary
		if strings.TrimSpace(line) != "" {
			// Clear the progress line and print the authentic output line
			fmt.Printf("\r\033[K  ➡️  %s\n", line)
		}
	}

	// Wait for the command to clean up and exit safely
	if err := cmd.Wait(); err != nil {
		fmt.Printf("\r\033[K❌ [%s] Command failed with error: %v\n", stepName, err)
		log.Fatalf("Execution halted.")
	}
	
	// Print a clean success mark clearing the progress text
	fmt.Printf("\r\033[K  ✅ Done!\n")
}
