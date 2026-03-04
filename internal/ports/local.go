/*
 *
 *  * Copyright 2026 P9 Labs
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *
 */

package ports

import "fmt"

type ListeningPort struct {
	Protocol string // "tcp"
	IP       string // "0.0.0.0", "127.0.0.1", etc.
	Port     int    // Port number
}

type ListeningPorts []ListeningPort

func (tcp ListeningPorts) Deduplication() []ListeningPort {
	result := make(map[string]ListeningPort)
	var ports []ListeningPort
	for _, p := range tcp {
		key := fmt.Sprintf("%s:%d", p.IP, p.Port)
		result[key] = p
	}
	for _, v := range result {
		ports = append(ports, v)
	}
	return ports
}

type BoundUDPPort struct {
	Protocol string // "udp"
	IP       string // "0.0.0.0", "127.0.0.1", etc.
	Port     int    // Port number
}

type BoundUDPPorts []BoundUDPPort

func (udp BoundUDPPorts) Deduplication() []BoundUDPPort {
	result := make(map[string]BoundUDPPort)
	var ports []BoundUDPPort
	for _, p := range udp {
		key := fmt.Sprintf("%s:%d", p.IP, p.Port)
		result[key] = p
	}
	for _, v := range result {
		ports = append(ports, v)
	}
	return ports
}

// GetListeningPorts returns all listening ports on the system
// Implementation is OS-specific (see local_linux.go, local_darwin.go)
