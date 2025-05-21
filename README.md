# REQD — Request Dispatcher

**reqd** is a fast, flexible HTTP POST request dispatcher written in Go. It reads JSON request bodies from a file and sends them concurrently to a target URL, making it useful for testing APIs or replaying large datasets of requests.

## Features

- ⚡ High-concurrency with configurable worker pool
- 🧠 Smart configuration: prioritizes CLI flags, then config file, then interactive prompts
- 📁 Reads request data from JSON file
- 🔐 Supports Bearer token authentication (more coming)
- 📊 Progress bar display for visibility into request processing
- 🧼 Graceful shutdown with context cancellation

## Installation

```sh
go install github.com/TRemigi/reqd
```

Make sure `$GOPATH/bin` or `$HOME/go/bin` is in your `PATH`.

## Usage

```sh
reqd -w 10 -u https://example.com/api -f ./input.json -t YOUR_TOKEN
```

You can also configure defaults in a config file at `~/.requester.conf`:

```ini
worker_count = 10
url = https://example.com/api
input_file = ./input.json
auth_token = YOUR_TOKEN
```

If any value is missing from the CLI args and config file, you’ll be prompted for it interactively.

## Flags

| Flag        | Description                        |
|-------------|------------------------------------|
| `-w`        | Number of concurrent workers       |
| `-u`        | Target URL                         |
| `-f`        | Path to JSON input file            |
| `-t`        | Bearer auth token                  |
| `-h`        | Display help message               |

## Input Format

Your input file should be a list of JSON objects:

```json
[
  {"name": "Alice", "email": "alice@example.com"},
  {"name": "Bob", "email": "bob@example.com"}
]
```

## License

MIT

## TODO

- [ ] Support custom token schemes
- [ ] Support all HTTP methods
- [ ] Retry logic with backoff
- [ ] Rate limiting
- [ ] Logging and metrics
- [ ] Make auth token optional
