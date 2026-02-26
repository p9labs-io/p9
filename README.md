# p9 - Network Debugging CLI Tool

> Network and Lookup "Swiss-knife"

A fast, lightweight network debugging tool written in Go for troubleshooting network connectivity, DNS lookups, and port scanning.

## ⚠️ Legal Disclaimer

**This tool is for authorized network diagnostic and educational purposes only.** P9 Labs is not responsible for any misuse or illegal activities. By using this tool, you agree to comply with all local and international laws. Do not use this tool to scan or probe networks or systems without explicit authorization.

## Features

### ✅ Implemented

#### Remote Port Checking (`-r`)
Check if remote TCP ports are open with customizable timeouts.

```bash
# Check if port is open
p9 -r google.com:443

# Check with custom timeout
p9 -r 192.168.1.1:22 -t 5s
```

**Features:**
- TCP port connectivity testing
- Colored output with status indicators
- Detailed error messages (connection refused, DNS errors, timeouts)
- Customizable timeout via `-t` flag

---

#### Local Open Ports (`-l`)
List all listening ports and services on your local machine.

```bash
# List all open ports
p9 -l
```

**Features:**
- Shows port number, protocol, and listening address
- TCP ports only (UDP support planned)
- Cross-platform support (Linux, macOS)
- IPv4 support

---

#### Domain Lookup (`-d`)
Query domain registration information using RDAP or WHOIS protocols.

```bash
# Auto mode: tries RDAP first, falls back to WHOIS
p9 -d example.com

# Force RDAP lookup
p9 -d example.com --rdap

# Force WHOIS lookup
p9 -d example.com --whois
```

**How it works:**
- By default, p9 tries RDAP first. If the TLD doesn't support RDAP, it automatically falls back to WHOIS.
- Use `--rdap` or `--whois` to force a specific protocol.
- Some ccTLDs (e.g. `.az`) have no public WHOIS or RDAP server — p9 will inform you and suggest alternatives where available.

##### RDAP Output

JSON output, easy to pipe to `jq`:

```bash
p9 -d example.com | jq '.events'
p9 -d example.com | jq '.nameservers'
p9 -d example.org | jq '.events[] | select(.eventAction=="expiration")'
p9 -d example.org | jq '.nameservers[].ldhName'
```

**Supported output fields:**
- `ldhName` - Domain name
- `events` - Registration, expiration, last changed dates
- `nameservers` - DNS servers
- `status` - Domain status codes
- `secureDNS` - DNSSEC configuration

##### WHOIS Output

Plain text output with key registration fields:
- Domain name and status
- Creation, expiry, and updated dates
- Registrar information
- Registrant, Admin, Tech contacts
- Name servers
- DNSSEC

**RDAP Configuration Cache:**

The tool downloads and caches RDAP server configurations from IANA. The config file location varies by platform:
- **Linux:** `~/.config/.p9/rdap_config`
- **macOS:** `~/Library/Application Support/.p9/rdap_config`
- **Windows:** `%APPDATA%\.p9\rdap_config`

This file is automatically refreshed every 30 days.

---

### 🚧 Planned Features

- [ ] **UDP port checking** - Test UDP port connectivity (part of `-r` functionality)
- [ ] **UDP port listing** - Show UDP listening ports in local scan (`-l`)
- [ ] **IPv6 support** - Local port checking for IPv6 addresses
- [ ] **IP geolocation** - Lookup IP address location and ASN information
- [ ] **Traceroute** - Network path visualization
- [ ] **DNS record lookup** - Query A, AAAA, MX, TXT records

---

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/p9labs-io/p9.git
cd p9

# Build
go build -o p9 ./cmd/

# (Optional) Install to system path
go install ./cmd/
```

### Binary Release

Download pre-built binaries from the [releases page](https://github.com/p9labs-io/p9/releases).

---

## Usage

```
p9 [flags]

Flags:
  -r string
        Check remote port (host:port)
  -t duration
        Override default timeout for port checks (default 3s)
  -l
        List local open ports
  -d string
        Domain lookup (RDAP with WHOIS fallback)
  --rdap
        Force RDAP lookup (use with -d)
  --whois
        Force WHOIS lookup (use with -d)
```

### Examples

```bash
# Check if SSH is open
p9 -r example.com:22

# Check with 10 second timeout
p9 -r slow-server.com:80 -t 10s

# List local listening services
p9 -l

# Domain lookup (auto: RDAP → WHOIS fallback)
p9 -d google.com

# Force WHOIS
p9 -d google.com --whois

# Force RDAP
p9 -d google.com --rdap

# Extract specific RDAP fields
p9 -d example.org | jq '.events[] | select(.eventAction=="expiration")'
p9 -d example.org | jq '.nameservers[].ldhName'
```

---

## Project Structure

```
p9/
├── cmd/
│   └── main.go           # CLI entry point
├── internal/
│   ├── cli/              # Output formatting
│   │   └── output.go
│   ├── dns/              # DNS/RDAP/WHOIS lookup
│   │   ├── rdap.go
│   │   └── whois.go
│   └── ports/            # Port scanning
│       ├── local.go
│       ├── local_darwin.go
│       ├── local_linux.go
│       └── remote.go
├── go.mod
└── README.md
```

---

## Requirements

- Go 1.21 or later (for building from source)
- No external dependencies (uses only Go standard library)

## Platform Support

- ✅ Linux
- ✅ macOS
- ⚠️ Windows (not fully tested)

---

## Contributing

Contributions are welcome! This is an educational project to learn Go and network programming.

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute.

## License

Copyright 2026 P9 Labs

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) file for details.

## Acknowledgments

- RDAP server routing based on [IANA RDAP Bootstrap Registry](https://data.iana.org/rdap/)
- WHOIS server routing based on [IANA WHOIS service](https://www.iana.org/whois)
- Inspired by traditional network debugging tools like `nc`, `nmap`, and `whois`
