/*
 * Copyright 2026 P9 Labs
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 */

package main

import (
	"flag"
	"github.com/p9labs-io/p9/internal/cli"
	"github.com/p9labs-io/p9/internal/ports"
	"time"
)

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
	//localFlag := flag.Bool("l", false, "List local open ports")
	//domainFlag := flag.String("d", "", "Domain/IP lookup")
	flag.Parse()

	switch {
	case *remoteFlag != "":
		result := ports.CheckPortTCP(*remoteFlag, *timeoutFlag)
		cli.PrintPortCheckResult(result)
	default:
		cli.PrintUsage()
		flag.PrintDefaults()
	}
}
