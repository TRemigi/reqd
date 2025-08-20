# REQD — Request Dispatcher

**reqd** is a fast, flexible HTTP request dispatcher written in Go. It reads JSON request data from a file and sends requests concurrently to a target URL using any HTTP method, making it useful for testing APIs or replaying large datasets of requests to test API versions.

## Features

- ⚡ High-concurrency with configurable worker pool
- 🧠 Smart configuration: prioritizes CLI flags, then config file, then interactive prompts
- 📁 Reads request data from a JSON file
- 🔐 Supports custom token schemes (e.g. Bearer)
- 📊 Progress bar display for visibility into request processing
- 🧼 Graceful shutdown with context cancellation

## Installation

```sh
go install github.com/TRemigi/reqd@latest
```

Make sure `$GOPATH/bin` or `$HOME/go/bin` is in your `PATH`.

## Usage

```sh
reqd -d ./data.json -rm POST -s Bearer -t YOUR_TOKEN -u https://example.com/api -w 10
```

You can also define defaults in a config file at `~/.reqd.conf`:

```ini
data_file = ./data.json
failure_log = falied.log
method = POST
success_log = succeeded.log
token_scheme = Bearer
token_value = YOUR_TOKEN
url = https://example.com/api
worker_count = 10
```

If any required values are missing from CLI flags and the config file, `reqd` will prompt you interactively.

## Flags

| Flag  | Description                                                                     |
| ----- | ------------------------------------------------------------------------------- |
| `-d`  | Path to JSON file containing an array of request data objects                   |
| `-rm`  | HTTP request method to use (`POST`, `GET`, `PUT`, `DELETE`, etc.)                       |
| `-lf` | Name of log file for failed requests (no failure logging if not provided)       |
| `-ls` | Name of log file for unsuccessful requests (no success logging if not provided) |
| `-s`  | Auth token scheme (e.g. `Bearer`)                                               |
| `-t`  | Auth token value                                                                |
| `-u`  | Target URL                                                                      |
| `-w`  | Number of concurrent worker goroutines to dispatch requests                     |
| `-h`  | Show help message and exit                                                      |

## Input Format

The data file should be a JSON array of objects. Each object represents one request's parameters:

```json
[
  { "name": "Alice", "email": "alice@example.com" },
  { "name": "Bob", "email": "bob@example.com" }
]
```

These are sent as request bodies for `POST`, query parameters for `GET`, etc., depending on the method.

## Configuration Precedence

Values are resolved in the following order:

1. Command-line flags
2. `~/.reqd.conf` config file
3. Interactive prompt

## Reporting

REQD supports optional logging of both **successful** and **failed** requests. If specified, results are written to the configured files in a simple human-readable format containing the full HTTP request and response.

To enable reporting, provide log file paths either via command-line flags or your config file:

### Flags

- `-ls <filename>` – log file for **successful** requests
- `-lf <filename>` – log file for **failed** requests

### Config file keys

```ini
success_log = ./success.rpt
failure_log = ./failures.rpt
```

Each log entry includes:

- The full HTTP request (method, URL, headers, body)
- The full HTTP response (status line, headers, body)

This is useful for debugging failed interactions or verifying behavior of successful ones.

> ℹ️ Log files will only be written if the respective option is configured. If neither is set, logging is skipped entirely.

## License

MIT

## TODO

- [x] Support custom token schemes
- [x] Support all HTTP methods
- [ ] Support concurrent execution of multiple config files
- [ ] Support combining request data from multiple data files
- [ ] Retry logic with backoff
- [ ] Rate limiting
- [ ] Logging and metrics
- [ ] Make auth token optional
