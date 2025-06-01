package main

import (
	"context"
	"flag"
	"fmt"
	"sync"

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
		flagFailureLog  = flag.String("lf", "", "Failure log file name")
		flagHelp        = flag.Bool("h", false, "Show help message")
		flagMethod      = flag.String("m", "", "Request method")
		flagSuccessLog  = flag.String("ls", "", "Success log file name")
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
		FailureLog:  *flagFailureLog,
		Method:      *flagMethod,
		SuccessLog:  *flagSuccessLog,
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

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	results := rex.GetReqd(c, jobs, ctx, cancel, &wg)
	reporting.ProcessResults(c, results, bar, &wg)

	wg.Wait()
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
