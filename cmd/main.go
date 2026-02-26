/*
 *
 *  * Copyright 2026 P9 Labs
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *
 */

package main

import (
	"flag"
	"fmt"
	"github.com/p9labs-io/p9/internal/cli"
	"github.com/p9labs-io/p9/internal/dns"
	"github.com/p9labs-io/p9/internal/ports"
	"time"
)

func runWhoisLookup(domain string, timeout time.Duration) {
	server, found := dns.GetWhoisServer(domain, timeout)
	if found {
		fmt.Print(dns.WhoisLookup(server, domain, timeout))
	} else {
		fmt.Println(server)
	}
}

func main() {
	// TODO: Define flags here
	/*
		1. Open ports and listening host - "p9 -l",
		2. Check remote port  - check remote tcp and udp ports "p9 -r 10.10.1.1:8080"
		3. domain whois and ip lookup - "p9 -d example.com"
	*/

	// TODO: Parse flags

	// TODO: Check which flag is set and route to the right function

	// For now, just print which operation was requested
	remoteFlag := flag.String("r", "", "Check remote port (host:port)")
	timeoutFlag := flag.Duration("t", 3*time.Second, "Override default timeout (e.g. -t 5s, -t 60s)")
	localFlag := flag.Bool("l", false, "List local open ports")
	domainFlag := flag.String("d", "", "Domain/IP lookup")
	rdapFlag := flag.Bool("rdap", false, "rdap Lookup")
	whoisFlag := flag.Bool("whois", false, "Whois Lookup")
	flag.Parse()

	switch {
	case *remoteFlag != "":
		result := ports.CheckPortTCP(*remoteFlag, *timeoutFlag)
		cli.PrintPortCheckResult(result)
	case *localFlag:
		result, err := ports.GetListeningPorts()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		cli.PrintListeningPorts(result)
	case *domainFlag != "":
		switch {
		case *whoisFlag:
			runWhoisLookup(*domainFlag, *timeoutFlag)
		case *rdapFlag:
			result, _ := dns.RdapResult(*domainFlag)
			fmt.Print(result)
		default:
			fmt.Println("Trying RDAP Servers...")
			result, supported := dns.RdapResult(*domainFlag)
			if supported {
				fmt.Print(result)
			} else {
				fmt.Println("Trying WHOIS Servers...")
				runWhoisLookup(*domainFlag, *timeoutFlag)
			}
		}
	default:
		cli.PrintUsage()
		flag.PrintDefaults()
	}
}
