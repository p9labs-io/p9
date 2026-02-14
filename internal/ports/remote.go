/*
 * Copyright 2026 P9 Labs
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 */

package ports

import (
	"errors"
	"net"
	"syscall"
	"time"
)

type PortCheckResult struct {
	Address   string
	IsOpen    bool
	Error     error
	ErrorType string // "timeout", "refused", "dns", "invalid_address", "other"
}

func CheckPortTCP(address string, timeout time.Duration) PortCheckResult {

	result := PortCheckResult{
		Address: address,
	}

	conn, err := net.DialTimeout("tcp", address, timeout)

	if err != nil {
		result.Error = err
		result.ErrorType = categorizeError(err)
		return result
	}
	defer conn.Close()

	result.IsOpen = true
	return result
}

func categorizeError(err error) string {
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		var netErr net.Error
		if errors.As(opErr.Err, &netErr) && netErr.Timeout() {
			return "timeout"
		}
		if errors.Is(opErr.Err, syscall.ECONNREFUSED) {
			return "refused"
		}
		var dnsErr *net.DNSError
		if errors.As(opErr.Err, &dnsErr) {
			return "dns"
		}
	}
	var addrErr *net.AddrError
	if errors.As(err, &addrErr) {
		return "invalid_address"
	}
	return "other"
}
