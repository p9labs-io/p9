# P9 - Network Diagnostic Tool

Network and Lookup "Swiss-knife"

## Features
- ✅ Check remote TCP ports
- 🚧 List local listening ports (in progress)
- 🚧 DNS lookup and WHOIS (planned)

## Installation

### From Source
````bash
go build -o p9 ./cmd/
````

## Usage

### Check Remote Port
````bash
./p9 -r google.com:443
````

**Output:**
````
✅ Port google.com:443 is OPEN
````

### List Local Ports
````bash
./p9 -l
````

### DNS Lookup
````bash
./p9 -d whois.com
````

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md)

## License

Apache 2.0 - See [LICENSE](LICENSE)

## Legal Disclaimer
This tool is for authorized network diagnostic and educational purposes only. P9 Labs is not responsible for any misuse or illegal activities. By using this tool, you agree to comply with all local and international laws.