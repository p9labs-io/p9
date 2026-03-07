/*
 *
 *  * Copyright 2026 P9 Labs
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *
 */

package ports

import (
	"fmt"
	"sort"
)

type ListeningPort struct {
	Command  string // "nginx"
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
	sort.Slice(ports, func(i, j int) bool { return ports[i].Port < ports[j].Port })
	return ports
}

type BoundUDPPort struct {
	Command  string // "Spotify"
	Protocol string // "udp"
	IP       string // "0.0.0.0", "127.0.0.1", etc.
	Port     int    // Port number
}

type BoundUDPPorts []BoundUDPPort

type LocalPorts interface {
	GetCommand() string
	GetProtocol() string
	GetIP() string
	GetPort() int
}

func (TCP ListeningPort) GetCommand() string {
	return TCP.Command
}

func (UDP BoundUDPPort) GetCommand() string {
	return UDP.Command
}

func (TCP ListeningPort) GetProtocol() string {
	return TCP.Protocol
}

func (UDP BoundUDPPort) GetProtocol() string {
	return UDP.Protocol
}

func (TCP ListeningPort) GetIP() string {
	return TCP.IP
}

func (UDP BoundUDPPort) GetIP() string {
	return UDP.IP
}

func (TCP ListeningPort) GetPort() int {
	return TCP.Port
}

func (UDP BoundUDPPort) GetPort() int {
	return UDP.Port
}

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
	sort.Slice(ports, func(i, j int) bool { return ports[i].Port < ports[j].Port })
	return ports
}
