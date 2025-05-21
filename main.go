package main

import (
	"flag"
	"fmt"

	"github.com/TRemigi/reqd/help"
	"github.com/TRemigi/reqd/reqconfig"
	"github.com/TRemigi/reqd/reqs"
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

	color.Green("___REQUESTER___")

	runConfig := reqconfig.GetConfig(flagWorkerCount, flagUrl, flagInputFile, flagAuthToken)
	reqconfig.PrintConfig(runConfig)

	reqBodies := rex.BodiesFromFile(runConfig.InputFile)
	bar := progressbar.New(len(reqBodies))
	jobs := createJobs(reqBodies)

	rex.GetReqd(runConfig, jobs, bar)
}

func createJobs(reqBodies []map[string]any) chan map[string]any {
	jobs := make(chan map[string]any, len(reqBodies))
	for _, item := range reqBodies {
		jobs <- item
	}
	close(jobs)
	return jobs
}
