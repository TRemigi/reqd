package reqconfig

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/TRemigi/reqd/pathutil"
	"github.com/fatih/color"
)

type RequestConfig struct {
	WorkerCount int
	Url         string
	InputFile   string
	AuthToken   string
}

func GetConfig(flagWorkerCount *int, flagUrl *string, flagInputFile *string, flagAuthToken *string) RequestConfig {
	fileConfig := configFromFile()
	argConfig := RequestConfig{
		WorkerCount: *flagWorkerCount,
		Url:         *flagUrl,
		InputFile:   *flagInputFile,
		AuthToken:   *flagAuthToken,
	}
	runConfig := runConfig(argConfig, *fileConfig)
	return runConfig
}

func PrintConfig(config RequestConfig) {
	color.Blue("Sending requests with configuration:")
	fmt.Fprintln(os.Stdout, "worker_count:", config.WorkerCount)
	fmt.Fprintln(os.Stdout, "url:", config.Url)
	fmt.Fprintln(os.Stdout, "input_file:", config.InputFile)
	fmt.Fprintln(os.Stdout, "auth_token:", config.AuthToken)
	fmt.Println()
}


func runConfig(argsConfig RequestConfig, fileConfig RequestConfig) RequestConfig {
	var workerCount int
	var url, inputFile, authToken string

	if argsConfig.WorkerCount != 0 {
		workerCount = argsConfig.WorkerCount
	} else if fileConfig.WorkerCount != 0 {
		workerCount = fileConfig.WorkerCount
	}

	if argsConfig.Url != "" {
		url = argsConfig.Url
	} else if fileConfig.Url != "" {
		url = fileConfig.Url
	}

	if argsConfig.InputFile != "" {
		inputFile = argsConfig.InputFile
	} else if fileConfig.InputFile != "" {
		inputFile = fileConfig.InputFile
	}

	if argsConfig.AuthToken != "" {
		authToken = argsConfig.AuthToken
	} else if fileConfig.AuthToken != "" {
		authToken = fileConfig.AuthToken
	}

	promptForMissingArgs(&workerCount, &url, &inputFile, &authToken)

	return RequestConfig{
		WorkerCount: workerCount, Url: url, InputFile: inputFile, AuthToken: authToken,
	}
}

func configFromFile() *RequestConfig {
	config := make(map[string]string)
	filePath := pathutil.ExpandPath("~/.requester.conf")
	file, err := os.Open(filePath)
	if err != nil {
		return nil // silently ignore missing file
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			config[key] = value
		}
	}
	WorkerCount, _ := strconv.Atoi(config["worker_count"])
	return &RequestConfig{
		WorkerCount: WorkerCount,
		Url:         config["url"],
		InputFile:   config["input_file"],
		AuthToken:   config["auth_token"],
	}
}

func promptForMissingArgs(WorkerCount *int, url *string, inputFile *string, authToken *string) {
	scanner := bufio.NewScanner(os.Stdin)
	if *WorkerCount == 0 {
		fmt.Print("Worker count: ")
		scanner.Scan()
		count, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		*WorkerCount = count
	}
	if *url == "" {
		fmt.Print("Enter request url: ")
		scanner.Scan()
		*url = scanner.Text()
	}
	if *inputFile == "" {
		fmt.Print("Enter request payload JSON file path: ")
		scanner.Scan()
		*inputFile = scanner.Text()
	}
	if *authToken == "" {
		fmt.Print("Enter auth token: ")
		scanner.Scan()
		*authToken = scanner.Text()
	}
}
