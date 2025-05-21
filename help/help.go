package help

func Help() string {
	return `___REQUESTER___

Usage:
  requester [worker_count] [url] [input_file] [auth_token]

Description:
  requester is a concurrent HTTP POST requester that reads JSON request bodies from a file and sends them to a target url.
  Useful for testing APIs or replaying request data at scale.

Argument Precedence:
  Arguments are resolved in the following order of priority:
    1. Command-line arguments
    2. Configuration file: ~/.requester.conf
    3. Interactive prompt (if any argument is missing)

Arguments:
  -w     Number of concurrent worker goroutines to send requests
  -u     Target url for POST requests
  -f     Path to a JSON file containing an array of request objects
  -t     Bearer token for Authorization header
	-h     Show this help message

Configuration File (~/.requester.conf):
  You may define default values in a simple key=value format:
    worker_count = 10
    url = http://localhost:8080/endpoint
    input_file = ~/Downloads/requests.json
    auth_token = abc123

Example:
  requester 5 https://api.example.com/data ./reqs.json supersecrettoken
  requester  # Uses values from ~/.requester.conf and prompts for missing ones

Help:
  -h   Show this help message and exit
`
}
