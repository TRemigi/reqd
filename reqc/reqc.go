package reqc

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
	Mode        string
	SuccessLog  string
	Token       string
	TokenScheme string
	Url         string
	WorkerCount int
	Headers     []Header
}

type Header struct {
	Name  string
	Value string
}

func Get(f RequestConfig, p string) RequestConfig {
	fileConfig := configFromFile(p)
	mergedConfig := mergedConfig(f, *fileConfig)
	promptForMissingRequiredArgs(&mergedConfig)
	PromptForAdditionalHeaders(&mergedConfig)
	return mergedConfig
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

func DataFromJSONFile(inputFile string) []map[string]any {
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
		Mode:        config["mode"],
		SuccessLog:  config["success_log"],
		Token:       config["token_value"],
		TokenScheme: config["token_scheme"],
		Url:         config["url"],
		WorkerCount: WorkerCount,
	}
}

func mergedConfig(argsConfig RequestConfig, fileConfig RequestConfig) RequestConfig {
	return RequestConfig{
		DataFile:    pick(argsConfig.DataFile, fileConfig.DataFile, ""),
		FailureLog:  pick(argsConfig.FailureLog, fileConfig.FailureLog, ""),
		Method:      pick(argsConfig.Method, fileConfig.Method, ""),
		Mode:        pick(argsConfig.Mode, fileConfig.Mode, ""),
		SuccessLog:  pick(argsConfig.SuccessLog, fileConfig.SuccessLog, ""),
		Token:       pick(argsConfig.Token, fileConfig.Token, ""),
		TokenScheme: pick(argsConfig.TokenScheme, fileConfig.TokenScheme, ""),
		Url:         pick(argsConfig.Url, fileConfig.Url, ""),
		WorkerCount: pick(argsConfig.WorkerCount, fileConfig.WorkerCount, 0),
	}
}

func pick[T comparable](primary T, fallback T, zero T) T {
	if primary != zero {
		return primary
	}
	return fallback
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

func promptForMissingRequiredArgs(reqconfig *RequestConfig) {
	scanner := bufio.NewScanner(os.Stdin)
	if reqconfig.WorkerCount == 0 {
		fmt.Print("Worker count: ")
		scanner.Scan()
		count, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		reqconfig.WorkerCount = count
	}
	if reqconfig.Url == "" {
		fmt.Print("Url: ")
		scanner.Scan()
		reqconfig.Url = scanner.Text()
	}
	if reqconfig.DataFile == "" {
		fmt.Print("JSON file path: ")
		scanner.Scan()
		reqconfig.DataFile = scanner.Text()
	}
	if reqconfig.Method == "" {
		fmt.Print("Method (post, get, put, delete): ")
		scanner.Scan()
		reqconfig.Method = scanner.Text()
	}
}

func PromptForAdditionalHeaders(reqconfig *RequestConfig) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Add additional headers? (y/N): ")
	scanner.Scan()
	addHeaders := scanner.Text()
	if addHeaders == "y" || addHeaders == "Y" {
		fmt.Print("Provide colon-separated header name and value: ")
		scanner.Scan()
		header := scanner.Text()
		parts := strings.SplitN(header, ":", 2)

		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			reqconfig.Headers = append(reqconfig.Headers, Header{Name: key, Value: value})
		}
	}
}
