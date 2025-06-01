package help

func Help() string {
	return `🦖 REQD - Concurrent HTTP Request Dispatcher

Usage:
  reqd [flags]

Description:
  REQD is a concurrent HTTP request dispatcher that reads request data from a JSON file
  and sends them to a target URL using the specified HTTP method.
  Useful for testing APIs or replaying request data at scale.

Flag Precedence:
  Arguments are resolved in the following order of priority:
    1. Command-line flags
    2. Configuration file: ~/.reqd.conf
    3. Interactive prompt (for missing required values)

Flags:
	-c     Path to config file (defaults to ~/.reqd.conf)
  -d     Path to a JSON file containing an array of request data objects
  -m     Request method (post, get, put, delete)
	-lf    Name of log file for failed requests
	-ls    Name of log file for successful requests
  -s     Auth token scheme (e.g., Bearer)
  -t     Auth token value
  -u     Target URL
  -w     Number of concurrent worker goroutines to send requests
  -h     Show this help message and exit

Configuration File (~/.reqd.conf):
  You may define default values using a simple key = value format:
    data_file = ~/Downloads/requests.json
		failure_log = falied.log
    method = POST
		success_log = succeeded.log
    token_scheme = Bearer
    token_value = abc123
    url = http://localhost:8080/endpoint
    worker_count = 8

Examples:
  reqd -d ./reqs.json -s Bearer -t supersecrettoken -u https://api.example.com/data -w 8 -lf failed.log -ls succeeded.log
  reqd                # Uses values from ~/.reqd.conf and prompts for any missing flags
`
}
