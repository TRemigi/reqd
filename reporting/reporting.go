package reporting

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

func Progress(results <-chan bool, bar *progressbar.ProgressBar, nreqs int) {
	totalFailed := 0
	for succeeded := range results {
		if !succeeded {
			totalFailed++
		}
		bar.Describe(fmt.Sprintf("%d/%d failed |", totalFailed, nreqs))
		bar.Add(1)
	}
}

func CreateReportFile() *os.File {
	wd, wderr := os.Getwd()
	if wderr != nil {
		log.Fatal(wderr)
	}
	f, err := os.Create(wd + "/" + time.Now().Truncate(time.Second).Format("2006-01-02_15-04-05") + ".rpt")
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func WriteToReport(reqBody []byte, respBody []byte, rFile *os.File) {
	report := fmt.Sprintf("Request:\n%s\nResponse:\n%s\n\n", string(reqBody), string(respBody))
	_, ferr := rFile.Write([]byte(report))
	if ferr != nil {
		log.Fatal(ferr)
	}
}
