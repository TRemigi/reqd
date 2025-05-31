package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/TRemigi/reqd/help"
	"github.com/TRemigi/reqd/reporting"
	"github.com/TRemigi/reqd/reqconfig"
	"github.com/TRemigi/reqd/rex"
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

	flags := reqconfig.Flags{
		WorkerCount: *flagWorkerCount,
		Url:         *flagUrl,
		InputFile:   *flagInputFile,
		AuthToken:   *flagAuthToken,
	}
	c := reqconfig.GetWithPrint(flags)

	reqBodies := reqconfig.BodiesFromFile(c.InputFile)
	numReqs := len(reqBodies)
	bar := progressbar.New(numReqs)
	jobs := rex.CreateJobs(reqBodies)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rf := reporting.CreateReportFile()
	results := rex.GetReqd(c, jobs, rf, ctx, cancel)
	reporting.Progress(results, bar, numReqs)
}
