package rex

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

func GetReqd(config reqconfig.RequestConfig, jobs chan map[string]any, ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) <-chan reporting.Result {
	results := make(chan reporting.Result)
	client := &http.Client{}

	var internalWg sync.WaitGroup
	for range config.WorkerCount {
		internalWg.Add(1)
		wg.Add(1)
		go func() {
			defer internalWg.Done()
			defer wg.Done()
			for job := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
					makeRequest(ctx, cancel, job, config, client, results)
				}
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		internalWg.Wait()
		close(results)
	}()

	return results
}

func makeRequest(ctx context.Context, cancel context.CancelFunc, job map[string]any, config reqconfig.RequestConfig, client *http.Client, results chan<- reporting.Result) {
	reqBody, _ := json.Marshal(job)
	req, _ := http.NewRequestWithContext(ctx, strings.ToUpper(config.Method), config.Url, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if config.TokenScheme != "" {
		req.Header.Set("Authorization", config.TokenScheme+" "+config.Token)
	}

	reqBuf := copyForLogging(req, reqBody)

	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			color.Yellow("Skipped request due to canceled context: %s", config.Url)
		} else {
			color.Red("Request failed, canceling context: %s", err)
			cancel()
		}
		return
	}
	defer resp.Body.Close()

	var resBuf bytes.Buffer
	resp.Write(&resBuf)

	results <- reporting.Result{
		Req:   reqBuf.Bytes(),
		Res:   resBuf.Bytes(),
		IsFailure: resp.StatusCode != http.StatusOK,
	}
}

func copyForLogging(req *http.Request, reqBody []byte) *bytes.Buffer {
	reqBuf := new(bytes.Buffer)
	fmt.Fprintf(reqBuf, "%s %s HTTP/1.1\r\n", req.Method, req.URL.String())
	req.Header.Write(reqBuf)
	fmt.Fprintf(reqBuf, "\r\n")
	reqBuf.Write(reqBody)
	return reqBuf
}
