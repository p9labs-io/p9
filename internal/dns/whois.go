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
	"os"
	"strings"
)

func GetWhoisServer(fqdn string) (string, bool) {
	conn, err := net.Dial("tcp", "whois.iana.org:43")
	if err != nil {
		log.Printf("IANA WHOIS server unavailable, %v\n", err)
	}
	tld := ExtractTLD(fqdn)
	query := fmt.Sprintf("%s\r\n", tld)
	conn.Write([]byte(query))

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
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	if whoisValue == "" {
		noWhois := fmt.Sprintf("No whois server provided, please use information below: \n%s", remarksValue)
		return noWhois, false
	}
	return fmt.Sprintf("%s", whoisValue), true
}

func WhoisLookup(srv string, domain string) {
	address := fmt.Sprintf("%s:43", srv)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Printf("WHOIS server unavailable, %v\n", err)
	}
	query := fmt.Sprintf("%s\r\n", domain)
	conn.Write([]byte(query))

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Domain") {
			value := strings.SplitN(line, ":", 2)
			fmt.Println(strings.TrimSpace(value[1]))
		}
		if strings.HasPrefix(line, "Registrar") {
			value := strings.SplitN(line, ":", 2)
			fmt.Println(strings.TrimSpace(value[1]))
		}
		if strings.HasPrefix(line, "Registrant") {
			value := strings.SplitN(line, ":", 2)
			fmt.Println(strings.TrimSpace(value[1]))
		}
		if strings.HasPrefix(line, "Admin") {
			value := strings.SplitN(line, ":", 2)
			fmt.Println(strings.TrimSpace(value[1]))
		}
		if strings.HasPrefix(line, "Tech") {
			value := strings.SplitN(line, ":", 2)
			fmt.Println(strings.TrimSpace(value[1]))
		}
		if strings.HasPrefix(line, "Name Server:") {
			value := strings.SplitN(line, ":", 2)
			fmt.Println(strings.TrimSpace(value[1]))
		}
		if strings.HasPrefix(line, "DNSSEC") {
			value := strings.SplitN(line, ":", 2)
			fmt.Println(strings.TrimSpace(value[1]))
		}

		//fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

//
//func whoisParser(line string) (string, bool) {
//	if strings.HasPrefix(line, "whois") {
//		value := strings.SplitN(line, ":", 2)
//		result := strings.TrimSpace(value[1])
//		fmt.Printf("%v", result)
//		if result == "" {
//			return "No whois server specified", false
//		}
//		return result, true
//	}
//	return "whois field not found", false
//}
//
//func remarksParser(line string) string {
//	for strings.HasPrefix(line, "remarks") {
//		value := strings.SplitN(line, ":", 2)
//		result := strings.TrimSpace(value[1])
//		return result
//	}
//	return "remarks field not found"
//}
