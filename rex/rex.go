package rex

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/TRemigi/reqd/reporting"
	"github.com/TRemigi/reqd/reqconfig"
	"github.com/fatih/color"
)

func CreateJobs(reqBodies []map[string]any) chan map[string]any {
	jobs := make(chan map[string]any, len(reqBodies))
	for _, item := range reqBodies {
		jobs <- item
	}
	close(jobs)
	return jobs
}

func GetReqd(config reqconfig.RequestConfig, jobs chan map[string]any, rFile *os.File, ctx context.Context, cancel context.CancelFunc) <-chan bool {
	results := make(chan bool)
	client := &http.Client{}
	var wg sync.WaitGroup

	for range config.WorkerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
					makeRequest(ctx, cancel, job, config, client, rFile, results)
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func makeRequest(ctx context.Context, cancel context.CancelFunc, job map[string]any, config reqconfig.RequestConfig, client *http.Client, rFile *os.File, results chan<- bool) {
	reqBody, _ := json.Marshal(job)
	req, _ := http.NewRequestWithContext(ctx, strings.ToUpper(config.Method), config.Url, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if config.TokenScheme != "" {
		req.Header.Set("Authorization", config.TokenScheme+" "+config.Token)
	}

	resp, err := client.Do(req)
	if err != nil {
		color.Red("Error: %s", err)
		cancel()
	} else if resp.StatusCode != http.StatusOK {
		results <- false
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		reporting.WriteToReport(reqBody, respBody, rFile)
	} else {
		results <- true
	}
}
