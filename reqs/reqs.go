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

func GetReqd(config reqconfig.RequestConfig, jobs chan map[string]any, bar *progressbar.ProgressBar) {
	var wg sync.WaitGroup
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
					makeRequest(ctx, job, config.Url, config.AuthToken, client, cancel)
					bar.Add(1)
				}
			}
		}()
	}
	wg.Wait()
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
