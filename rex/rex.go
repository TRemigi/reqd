package rex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/TRemigi/reqd/pathutil"
	"github.com/TRemigi/reqd/reqconfig"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

func GetReqd(config reqconfig.RequestConfig, jobs chan map[string]any, bar *progressbar.ProgressBar, rFile *os.File) {
	nreqs := len(jobs)
	var wg sync.WaitGroup
	failures := failures()
	client := &http.Client{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	color.Blue("Running...")
	for range config.WorkerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
					makeRequest(ctx, job, config.Url, config.AuthToken, client, cancel, rFile, failures)
					bar.Add(1)
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(failures)
	}()

	totalFailed := 0
	for range failures {
			totalFailed++
			bar.Describe(fmt.Sprintf("%d/%d failed |", totalFailed, nreqs))
	}
}

func BodiesFromFile(inputFile string) []map[string]any {
	data, err := os.ReadFile(pathutil.ExpandPath(inputFile))
	if err != nil {
		log.Fatal(err)
	}

	var items []map[string]any
	if err := json.Unmarshal(data, &items); err != nil {
		log.Fatal(err)
	}
	return items
}

func failures() chan int {
	failures := make(chan int)
	return failures
}

func makeRequest(ctx context.Context, job map[string]any, url string, authToken string, client *http.Client, cancel context.CancelFunc, rFile *os.File, failures chan<- int) {
	reqBody, _ := json.Marshal(job)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		color.Red("Error: %s", err)
		cancel()
	} else if resp.StatusCode != http.StatusOK {
		failures <- 1

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		writeToReport(reqBody, respBody, rFile)
	}
}

func writeToReport(reqBody []byte, respBody []byte, rFile *os.File) {
	report := fmt.Sprintf("Request:\n%s\nResponse:\n%s\n\n", string(reqBody), string(respBody))
	_, ferr := rFile.Write([]byte(report))
	if ferr != nil {
		log.Fatal(ferr)
	}
}
