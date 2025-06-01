package reqconfig

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/TRemigi/reqd/pathutil"
)

type RequestConfig struct {
	DataFile    string
	FailureLog  string
	Method      string
	SuccessLog  string
	Token       string
	TokenScheme string
	Url         string
	WorkerCount int
}

func Get(f RequestConfig, p string) RequestConfig {
	fileConfig := configFromFile(p)
	return finalConfig(f, *fileConfig)
}

func Print(c RequestConfig) {
	fmt.Printf(" :: Data File    : %s\n", c.DataFile)
	fmt.Printf(" :: Method       : %s\n", strings.ToUpper(c.Method))
	fmt.Printf(" :: Token Scheme : %s\n", c.TokenScheme)
	fmt.Printf(" :: Token Value  : %s\n", redactToken(c.Token))
	fmt.Printf(" :: URL          : %s\n", c.Url)
	fmt.Printf(" :: Worker Count : %d\n", c.WorkerCount)
	fmt.Printf(" :: SuccessLog   : %s\n", c.SuccessLog)
	fmt.Printf(" :: FailureLog   : %s\n", c.FailureLog)
	fmt.Println("________________________________________________")
	fmt.Println()
}

func GetWithPrint(f RequestConfig, p string) RequestConfig {
	c := Get(f, p)
	Print(c)
	return c
}

func DataFromFile(inputFile string) []map[string]any {
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

func configFromFile(p string) *RequestConfig {
	config := make(map[string]string)
	var filePath string
	if p != "" {
		filePath = p
	} else {
		filePath = pathutil.ExpandPath("~/.reqd.conf")
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
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
		DataFile:    config["data_file"],
		FailureLog:  config["failure_log"],
		Method:      config["method"],
		SuccessLog:  config["success_log"],
		Token:       config["token_value"],
		TokenScheme: config["token_scheme"],
		Url:         config["url"],
		WorkerCount: WorkerCount,
	}
}

func finalConfig(argsConfig RequestConfig, fileConfig RequestConfig) RequestConfig {
	var dataFile, failureLog, method, successLog, tokenScheme, tokenValue, url string
	var workerCount int

	if argsConfig.DataFile != "" {
		dataFile = argsConfig.DataFile
	} else if fileConfig.DataFile != "" {
		dataFile = fileConfig.DataFile
	}

	if argsConfig.FailureLog != "" {
		failureLog = argsConfig.FailureLog
	} else if fileConfig.FailureLog != "" {
		failureLog = fileConfig.FailureLog
	}

	if argsConfig.Method != "" {
		method = argsConfig.Method
	} else if fileConfig.Method != "" {
		method = fileConfig.Method
	}

	if argsConfig.SuccessLog != "" {
		successLog = argsConfig.SuccessLog
	} else if fileConfig.SuccessLog != "" {
		successLog = fileConfig.SuccessLog
	}

	if argsConfig.Token != "" {
		tokenValue = argsConfig.Token
	} else if fileConfig.Token != "" {
		tokenValue = fileConfig.Token
	}

	if argsConfig.TokenScheme != "" {
		tokenScheme = argsConfig.TokenScheme
	} else if fileConfig.TokenScheme != "" {
		tokenScheme = fileConfig.TokenScheme
	}

	if argsConfig.Url != "" {
		url = argsConfig.Url
	} else if fileConfig.Url != "" {
		url = fileConfig.Url
	}

	if argsConfig.WorkerCount != 0 {
		workerCount = argsConfig.WorkerCount
	} else if fileConfig.WorkerCount != 0 {
		workerCount = fileConfig.WorkerCount
	}

	promptForMissingRequiredArgs(&workerCount, &method, &url, &dataFile)

	return RequestConfig{
		DataFile:    dataFile,
		FailureLog:  failureLog,
		Method:      method,
		SuccessLog:  successLog,
		Token:       tokenValue,
		TokenScheme: tokenScheme,
		Url:         url,
		WorkerCount: workerCount,
	}
}

func redactToken(token string) string {
	if len(token) == 0 {
		return ""
	}
	if len(token) <= 4 {
		return "****"
	}
	return "****" + token[len(token)-4:]
}

func promptForMissingRequiredArgs(workerCount *int, method *string, url *string, dataFile *string) {
	scanner := bufio.NewScanner(os.Stdin)
	if *workerCount == 0 {
		fmt.Print("Worker count: ")
		scanner.Scan()
		count, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		*workerCount = count
	}
	if *url == "" {
		fmt.Print("Url: ")
		scanner.Scan()
		*url = scanner.Text()
	}
	if *dataFile == "" {
		fmt.Print("JSON file path: ")
		scanner.Scan()
		*dataFile = scanner.Text()
	}
	if *method == "" {
		fmt.Print("Method (post, get, put, delete): ")
		scanner.Scan()
		*method = scanner.Text()
	}
}
