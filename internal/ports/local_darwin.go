//go:build darwin

package ports

// NOTE: Currently uses lsof command. Future optimization: use syscalls for better performance.

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func GetListeningPorts() ([]ListeningPort, error) {
	if _, err := exec.LookPath("lsof"); err != nil {
		return nil, fmt.Errorf("lsof command not found - this tool requires lsof on macOS")
	}

	cmd := exec.Command("lsof", "-iTCP", "-sTCP:LISTEN", "-n", "-P")

	output, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var ports []ListeningPort
	scanner := bufio.NewScanner(output)

	// Skip header line
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()

		fields := strings.Fields(line)

		if len(fields) < 9 {
			continue
		}

		nameField := fields[8]
		protoField := strings.ToLower(fields[7])

		splittedNameField := strings.SplitN(nameField, ":", 2)

		if splittedNameField[0] == "*" {
			splittedNameField[0] = "0.0.0.0"
		}

		ip := splittedNameField[0]
		// Convert port string to int
		p, err := strconv.Atoi(splittedNameField[1])
		if err != nil {
			log.Printf("Warning: Skipping malformed port in line: %s (error: %v)\n", nameField, err)
			continue
		}

		ports = append(ports, ListeningPort{Port: p, Protocol: protoField, IP: ip})
	}

	cmd.Wait()
	return ports, scanner.Err()
}
