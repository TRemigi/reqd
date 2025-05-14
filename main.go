package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/TRemigi/requester/path"
	"github.com/TRemigi/requester/reqconfig"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

func main() {
	var (
		flagWorkerCount = flag.Int("w", 0, "Number of concurrent workers")
		flagUrl         = flag.String("u", "", "Target url for POST requests")
		flagInputFile   = flag.String("f", "", "Path to JSON input file")
		flagAuthToken   = flag.String("t", "", "Bearer auth token")
		flagHelp        = flag.Bool("h", false, "Show help message")
	)
	flag.Parse()

	if *flagHelp != false {
		fmt.Println(getHelpMessage())
		return
	}

	color.Green("___REQUESTER___")

	// Gather args in priority: CLI > config > prompt
	runConfig := reqconfig.GetConfig(flagWorkerCount, flagUrl, flagInputFile, flagAuthToken)
	reqconfig.PrintConfig(runConfig)

	reqBodies := bodiesFromFile(runConfig.InputFile)
	bar := progressbar.New(len(reqBodies))
	jobs := createJobs(reqBodies)

	var wg sync.WaitGroup
	client := &http.Client{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	color.Blue("Running...")
	for range runConfig.WorkerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
					makeRequest(ctx, job, runConfig.Url, runConfig.AuthToken, client, cancel)
					bar.Add(1)
				}
			}
		}()
	}
	wg.Wait()
}

func bodiesFromFile(inputFile string) []map[string]any {
	data, err := os.ReadFile(path.ExpandPath(inputFile))
	if err != nil {
		log.Fatal(err)
	}

	var items []map[string]any
	if err := json.Unmarshal(data, &items); err != nil {
		log.Fatal(err)
	}
	return items
}

func createJobs(reqBodies []map[string]any) chan map[string]any {
	jobs := make(chan map[string]any, len(reqBodies))
	for _, item := range reqBodies {
		jobs <- item
	}
	close(jobs)
	return jobs
}

func makeRequest(ctx context.Context, job map[string]any, url string, authToken string, client *http.Client, cancel context.CancelFunc) {
	reqBody, _ := json.Marshal(job)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		color.Red("Error: %s", err)
		cancel()
	} else if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		fmt.Println()
		color.Yellow("Unsuccessful request:")
		fmt.Println(string(reqBody))
		color.Yellow("Response Body:")
		fmt.Println(string(respBody))
		fmt.Println()
	}
}

func getHelpMessage() string {
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
