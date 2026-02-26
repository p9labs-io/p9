/*
 *
 *  * Copyright 2026 P9 Labs
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *
 */

package dns

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func GetWhoisServer(fqdn string, timeout time.Duration) (string, bool) {
	//TODO: implement unknown tlds message

	conn, err := net.DialTimeout("tcp", "whois.iana.org:43", timeout)
	if err != nil {
		log.Fatalf("IANA WHOIS server unavailable, %v\n", err)
	}
	defer conn.Close()
	tld := ExtractTLD(fqdn)
	query := fmt.Sprintf("%s\r\n", tld)
	_, err = conn.Write([]byte(query))
	if err != nil {
		log.Printf("Error writing to connection:\n%v", err)
	}

	var whoisValue, remarksValue string
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "whois") {
			value := strings.SplitN(line, ":", 2)
			whoisValue = strings.TrimSpace(value[1])
		}
		if strings.HasPrefix(line, "remarks") {
			value := strings.SplitN(line, ":", 2)
			remarksValue = strings.TrimSpace(value[1])
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading response: %v", err)
	}
	if whoisValue == "" {
		noWhois := fmt.Sprintf("No whois server provided, please use information below: \n%s", remarksValue)
		//TODO: multiple remarks case
		return noWhois, false
	}
	return whoisValue, true
}

func WhoisLookup(srv string, domain string, timeout time.Duration) string {
	address := fmt.Sprintf("%s:43", srv)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		log.Fatalf("WHOIS server unavailable, %v\n", err)
	}
	defer conn.Close()
	query := fmt.Sprintf("%s\r\n", domain)
	_, err = conn.Write([]byte(query))
	if err != nil {
		log.Printf("Error writing to connection:\n%v", err)
	}

	var output []string
	prefix := []string{"Domain", "Creation Date", "Registry Expiry Date", "Updated Date", "Registrar ", "Registrant", "Admin", "Tech", "Name Server", "DNSSEC"}
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		for _, v := range prefix {
			if strings.HasPrefix(line, v) {
				output = append(output, line)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading standard input: %v", err)
	}
	return strings.Join(output, "\n")
}
