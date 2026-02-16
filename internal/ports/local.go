/*
 *
 *  * Copyright 2026 P9 Labs
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *
 */

package ports

type ListeningPort struct {
	Protocol string // "tcp"
	IP       string // "0.0.0.0", "127.0.0.1", etc.
	Port     int    // Port number
}

// GetListeningPorts returns all listening ports on the system
// Implementation is OS-specific (see local_linux.go, local_darwin.go)
