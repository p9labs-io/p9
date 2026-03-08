/*
 *
 *  * Copyright 2026 P9 Labs
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *
 */

//go:build darwin || linux

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

func parseNameField(nameField string) (ip string, port string) {
	splittedNameField := strings.Split(nameField, ":")

	ip = strings.Join(splittedNameField[:len(splittedNameField)-1], ":")
	ip = strings.Trim(ip, "[]")
	port = splittedNameField[len(splittedNameField)-1]
	if ip == "*" {
		ip = "0.0.0.0"
	}

	return ip, port
}

func GetListeningPorts() (ListeningPorts, error) {
	if _, err := exec.LookPath("lsof"); err != nil {
		return nil, fmt.Errorf("lsof not found — install with: brew install lsof / apt install lsof / yum install lsof\n")
	}

	cmd := exec.Command("lsof", "-iTCP", "-sTCP:LISTEN", "-n", "-P")

	output, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var ports ListeningPorts
	scanner := bufio.NewScanner(output)

	// Skip header line
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()

		fields := strings.Fields(line)

		if len(fields) < 9 {
			log.Printf("Something went wrong with lsof output fields")
			continue
		}

		nameField := fields[8]
		protoField := strings.ToLower(fields[7])
		commandField := strings.ToLower(fields[0])

		ip, port := parseNameField(nameField)
		// Convert port string to int
		p, err := strconv.Atoi(port)
		if err != nil {
			log.Printf("Warning: Skipping malformed port in line: %s (error: %v)\n", nameField, err)
			continue
		}

		ports = append(ports, ListeningPort{Command: commandField, Port: p, Protocol: protoField, IP: ip})
	}

	cmd.Wait()
	return ports, scanner.Err()
}

func GetBoundUDPPorts() (BoundUDPPorts, error) {
	if _, err := exec.LookPath("lsof"); err != nil {
		return nil, fmt.Errorf("lsof not found — install with: brew install lsof / apt install lsof / yum install lsof\n")
	}

	cmd := exec.Command("lsof", "-iUDP", "-n", "-P")

	output, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var ports []BoundUDPPort
	scanner := bufio.NewScanner(output)

	// Skip header line
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "->") {
			continue
		}

		fields := strings.Fields(line)

		if len(fields) < 9 {
			log.Printf("Something went wrong with lsof output fields")
			continue
		}

		nameField := fields[8]
		protoField := strings.ToLower(fields[7])
		commandField := strings.ToLower(fields[0])

		ip, port := parseNameField(nameField)
		// Convert port string to int
		if port != "*" {
			p, err := strconv.Atoi(port)
			if err != nil {
				log.Printf("Warning: Skipping malformed port in line: %s (error: %v)\n", nameField, err)
				continue
			}
			ports = append(ports, BoundUDPPort{Command: commandField, Port: p, Protocol: protoField, IP: ip})
		}
	}

	cmd.Wait()
	return ports, scanner.Err()
}
