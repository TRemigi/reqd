package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/TRemigi/reqd/help"
	"github.com/TRemigi/reqd/reporting"
	"github.com/TRemigi/reqd/reqconfig"
	"github.com/TRemigi/reqd/rex"
	"github.com/TRemigi/reqd/version"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

func main() {
	var (
		flagConfigFile  = flag.String("c", "", "Path to config file")
		flagDataFile    = flag.String("d", "", "Path to JSON data file")
		flagHelp        = flag.Bool("h", false, "Show help message")
		flagMethod      = flag.String("m", "", "Request method")
		flagToken       = flag.String("t", "", "Auth token value")
		flagTokenScheme = flag.String("s", "", "Auth token scheme")
		flagUrl         = flag.String("u", "", "Target url")
		flagWorkerCount = flag.Int("w", 0, "Number of concurrent workers")
	)
	flag.Parse()

	if *flagHelp != false {
		fmt.Println(help.Help())
		return
	}

	printHeader(version.Version)

	fConfig := reqconfig.RequestConfig{
		DataFile:    *flagDataFile,
		Method:      *flagMethod,
		Token:       *flagToken,
		TokenScheme: *flagTokenScheme,
		Url:         *flagUrl,
		WorkerCount: *flagWorkerCount,
	}
	c := reqconfig.GetWithPrint(fConfig, *flagConfigFile)

	reqData := reqconfig.DataFromFile(c.DataFile)
	numReqs := len(reqData)
	bar := progressbar.New(numReqs)
	jobs := rex.CreateJobs(reqData)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rf := reporting.CreateReportFile()
	results := rex.GetReqd(c, jobs, rf, ctx, cancel)
	reporting.Progress(results, bar, numReqs)
}

func printHeader(version string) {
	color.Green(`
  ██████╗ ███████╗ ██████╗ ██████╗ 
  ██╔══██╗██╔════╝██╔═══██╗██╔══██╗
  ██████╔╝█████╗  ██║   ██║██║  ██║
  ██╔══██╗██╔══╝  ██║▄▄ ██║██║  ██║
  ██║  ██║███████╗╚██████╔╝██████╔╝
  ╚═╝  ╚═╝╚══════╝ ╚══▀▀═╝ ╚═════╝ 
`)
	fmt.Printf("REQD — Request Dispatcher v%s\n", version)
	fmt.Println("________________________________________________")
	fmt.Println()
}
