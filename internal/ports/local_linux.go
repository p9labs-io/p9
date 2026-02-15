//go:build linux

package ports

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetListeningPorts() ([]ListeningPort, error) {
	var allPorts []ListeningPort

	tcp4Ports, err := parsePortFile("/proc/net/tcp", "tcp")
	if err != nil {
		return nil, err
	}
	allPorts = append(allPorts, tcp4Ports...)

	tcp6Ports, err := parsePortFile("/proc/net/tcp6", "tcp6")
	if err != nil {
		// IPv6 might not exist, that's okay
	} else {
		allPorts = append(allPorts, tcp6Ports...)
	}

	return allPorts, nil
}

func parsePortFile(filename string, protocol string) ([]ListeningPort, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ports []ListeningPort
	scanner := bufio.NewScanner(file)

	// Skip header
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 4 {
			continue
		}

		// Check if listening (state 0A)
		if fields[3] != "0A" {
			continue
		}

		// Parse local_address
		parts := strings.Split(fields[1], ":")
		if len(parts) != 2 {
			continue
		}

		hexIP := parts[0]
		hexPort := parts[1]

		// Decode port
		port, err := strconv.ParseInt(hexPort, 16, 32)
		if err != nil {
			continue
		}

		// Decode IP
		ip, err := decodeIPv4(hexIP)
		if err != nil {
			continue
		}

		ports = append(ports, ListeningPort{
			Protocol: protocol,
			IP:       ip,
			Port:     int(port),
		})
	}

	return ports, scanner.Err()
}

func decodeIPv4(hexIP string) (string, error) {
	bytes, err := hex.DecodeString(hexIP)
	if err != nil {
		return "", err
	}

	if len(bytes) != 4 {
		return "", fmt.Errorf("invalid IP length")
	}

	// Reverse bytes (little-endian)
	for i := 0; i < len(bytes)/2; i++ {
		j := len(bytes) - 1 - i
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}

	return fmt.Sprintf("%d.%d.%d.%d", bytes[0], bytes[1], bytes[2], bytes[3]), nil
}
