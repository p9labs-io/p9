/*
 *
 *  * Copyright 2026 P9 Labs
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *
 */

package cli

import (
	"fmt"
	"github.com/p9labs-io/p9/internal/dns"
	"github.com/p9labs-io/p9/internal/ports"
	"strings"
)

func PrintPortCheckResult(result ports.PortCheckResult) {
	if result.IsOpen {
		fmt.Printf("✅ \033[32mPort %s is OPEN\033[0m\n", result.Address)
		return
	}

	switch result.ErrorType {
	case "timeout":
		fmt.Printf("⏱️  \033[33mConnection to %s TIMED OUT\033[0m\n", result.Address)
	case "refused":
		fmt.Printf("🔴 \033[31mPort %s is CLOSED (connection refused)\033[0m\n", result.Address)
	case "dns":
		fmt.Printf("🔍 \033[33mDNS lookup failed: %v\033[0m\n", result.Error)
	case "invalid_address":
		fmt.Printf("⚠️  \033[33mInvalid address format: %v\033[0m\n", result.Error)
	default:
		fmt.Printf("❌ \033[31mCannot connect to %s: %v\033[0m\n", result.Address, result.Error)
	}
}

func PrintUsage() {
	fmt.Println("Debugging tool")
	fmt.Println("\nUsage:")
	fmt.Println("  p9 -r <host:port>  Check if remote port is open")
	fmt.Println("  p9 -l              List local open ports")
	fmt.Println("  p9 -d <domain>     Lookup domain/IP information")
}

func PrintListeningPorts(ports []ports.ListeningPort) {
	if len(ports) == 0 {
		fmt.Println("No listening ports found")
		return
	}

	fmt.Println("Listening ports:")
	for _, p := range ports {
		fmt.Printf("  %s %s:%d\n", strings.ToUpper(p.Protocol), p.IP, p.Port)
	}
}

func PrintWhoisResult(domain string) {
	server, found := dns.GetWhoisServer(domain)
	if found {
		dns.WhoisLookup(server, domain)
	} else {
		fmt.Println(server)
	}
}
