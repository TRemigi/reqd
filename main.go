package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TRemigi/reqd/help"
	"github.com/TRemigi/reqd/reqconfig"
	rex "github.com/TRemigi/reqd/reqs"
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
		fmt.Println(help.Help())
		return
	}

	color.Green("🦖 REQD")
	fmt.Println()

	runConfig := reqconfig.GetConfig(flagWorkerCount, flagUrl, flagInputFile, flagAuthToken)
	reqconfig.PrintConfig(runConfig)

	reqBodies := rex.BodiesFromFile(runConfig.InputFile)
	bar := progressbar.New(len(reqBodies))
	jobs := createJobs(reqBodies)

	rf := createReportFile()
	rex.GetReqd(runConfig, jobs, bar, rf)
}

func createJobs(reqBodies []map[string]any) chan map[string]any {
	jobs := make(chan map[string]any, len(reqBodies))
	for _, item := range reqBodies {
		jobs <- item
	}
	close(jobs)
	return jobs
}

func createReportFile() *os.File {
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
