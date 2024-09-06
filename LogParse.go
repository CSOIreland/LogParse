package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Result struct {
	LogFile       string
	Time          string
	Origin        string
	CorrelationId string
	LineNumber    int
	Error         string
	Frequency     int
}

var lineBreakRegExp = regexp.MustCompile(`\r\n`)
var results []Result = []Result{}

func main() {
	var directory = os.Args[1]
	if len(directory) == 0 {
		directory = "."
	}
	log.Println("Parsing log files in", directory)

	var files, err = os.ReadDir(directory)
	if err != nil {
		log.Fatalln(err)
	}

	var fileInfos = make([]os.FileInfo, 0)
	for _, file := range files {
		var fileInfo, err = file.Info()
		if err != nil {
			log.Fatalln(err)
		}
		fileInfos = append(fileInfos, fileInfo)
	}

	// Sort in ascending order
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].ModTime().After(fileInfos[j].ModTime())
	})
	var errors = make(map[string]string)
	for _, fileInfo := range fileInfos {
		readFile(directory+string(os.PathSeparator)+fileInfo.Name(), errors)
	}
	fmt.Println(fmt.Sprintf("Number of errors: %d", len(errors)))
	for k, v := range errors {
		parse(errors, k, v)
	}

	// Sort Results by frequency
	sort.Slice(results, func(i, j int) bool {
		return results[i].Frequency > results[j].Frequency
	})

	// Write results as JSON to file
	file, _ := json.MarshalIndent(results, "", " ")
	var now = time.Now().Format("2006-02-01")
	var outputFileName = "results" + now + ".json"
	var wd, _ = os.Getwd()
	_ = os.WriteFile(outputFileName, file, 0200)
	fmt.Println(fmt.Sprintf("Output Filename is %s", wd+string(os.PathSeparator)+outputFileName))
}

func readFile(name string, errors map[string]string) {
	fileContents, err := os.ReadFile(name)
	if err != nil {
		log.Fatalln(err)
	}

	fileLines := lineBreakRegExp.Split(string(fileContents), -1)

	fmt.Println(fmt.Sprintf("Filename is %s", name))
	var s strings.Builder
	for i, line := range fileLines {
		if strings.Contains(line, "ERROR") || strings.Contains(line, "FATAL") {
			s.WriteString(fmt.Sprintf("%s", line))
			errors[name+"$"+strconv.Itoa(i+1)] = s.String()
			s.Reset()
		}
	}
}

func parse(errors map[string]string, k string, v string) {
	var logFile = strings.Split(k, "$")[0]
	var lineNumber, _ = strconv.Atoi(strings.Split(k, "$")[1])
	var time = parseTime(v)
	var origin = parseOrigin(v)
	var correlationId = parseCorrelationId(v)
	var error = parseError(v)
	var frequency = 0

	if len(errors) == 1 {
		frequency = 1
	} else {
		for key, value := range errors {
			var test = strings.Split(value, "] ")[1]

			// Compare the first 50 characters of error and test and if they are the same, assume that the errors
			// refer to the same issue
			if len(error) < 50 || len(test) < 50 {
				continue
			}
			if error[:50] == test[:50] {
				frequency++
				delete(errors, key)
			}
		}
	}

	var result = Result{
		LogFile:       logFile,
		LineNumber:    lineNumber,
		Time:          time,
		Origin:        origin,
		CorrelationId: correlationId,
		Error:         error,
		Frequency:     frequency,
	}
	results = append(results, result)
}

func parseTime(error string) string {

	// Regular expression to match everything before the first square bracket
	re := regexp.MustCompile(`^(.*?)\s*\[`)

	// Find the match
	match := re.FindStringSubmatch(error)

	// Check if a match is found
	if len(match) > 1 {
		// Print the part before the first square bracket
		return match[1]
	}
	return ""
}

func parseOrigin(error string) string {

	// Regular expression to match the content inside the first square bracket
	re := regexp.MustCompile(`\[(.*?)\]`)

	// Find the first match
	match := re.FindStringSubmatch(error)

	// Check if a match is found
	if len(match) > 1 {
		// Print the content inside the first square bracket
		return match[1]
	}
	return ""
}

func parseCorrelationId(error string) string {

	// Regular expression to match the content inside square brackets
	re := regexp.MustCompile(`\[(.*?)\]`)

	// Find all matches
	matches := re.FindAllStringSubmatch(error, -1)

	// Check if there are at least two matches
	if len(matches) >= 2 {
		// Print the content inside the second square bracket
		return matches[1][1]
	}
	return ""
}

func parseError(error string) string {

	// Regular expression to match content within square brackets
	re := regexp.MustCompile(`\[(.*?)\]`)

	// Find all matches
	matches := re.FindAllStringSubmatchIndex(error, -1)

	// Check if there are at least three matches
	if len(matches) >= 3 {
		// Get the position right after the third match (after the closing bracket)
		errorIndex := matches[2][1]
		return error[errorIndex:]
	}
	return ""
}
