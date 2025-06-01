package reporting

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/TRemigi/reqd/reqconfig"
	"github.com/schollz/progressbar/v3"
)

type Result struct {
	Req   []byte
	Res   []byte
	IsFailure bool
}

func ProcessResults(c reqconfig.RequestConfig, results <-chan Result, bar *progressbar.ProgressBar, wg *sync.WaitGroup) {
	failureChan, successChan := make(chan Result, 500), make(chan Result, 500)
	if c.FailureLog != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			startLogger(failureChan, c.FailureLog)
		}()
	}

	if c.SuccessLog != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			startLogger(successChan, c.SuccessLog)
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		totalFailed := 0
		for result := range results {
			if result.IsFailure {
				totalFailed++
				if c.FailureLog != "" {
					failureChan <- result
				}
			} else {
				if c.SuccessLog != "" {
					successChan <- result
				}
			}
			progress(totalFailed, bar)
		}
		close(failureChan)
		close(successChan)
	}()
}

func startLogger(results <-chan Result, n string) {
	f := createReportFile(n)
	for result := range results {
		reportResult(result, f)
	}
}

func progress(nFailed int, bar *progressbar.ProgressBar) {
	bar.Describe(fmt.Sprintf("%d failed |", nFailed))
	bar.Add(1)
}

func createReportFile(name string) *os.File {
	wd, wderr := os.Getwd()
	if wderr != nil {
		log.Fatal(wderr)
	}
	f, err := os.Create(wd + "/" + name)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func reportResult(r Result, f *os.File) {
	req := fmt.Sprintf("Request:\n%s\n", string(r.Req))
	writeToReport(f, req)
	writeToReport(f, "\n")
	res := fmt.Sprintf("Response:\n%s\n", string(r.Res))
	writeToReport(f, res)
	writeToReport(f, fmt.Sprintf("==================================\n"))
	writeToReport(f, "\n")
}

func writeToReport(f *os.File, report string) {
	_, err := f.Write([]byte(report))
	if err != nil {
		log.Fatal(err)
	}
}
