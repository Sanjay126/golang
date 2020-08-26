package main

import (
	"encoding/csv"
	"flag"
	"github.com/gocarina/gocsv"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	var writer *csv.Writer
	var csvFile *os.File
	var users []*event.Payload
	fileArg := flag.String("file", "set.csv", "CSV file name")
	startIndex := flag.Int("start", 0, "Start Index")
	rlArg := flag.Int("rl", 10, "Rate Limit")
	concurrency := flag.Int("concurrency", 20, "Concurrent executors")
	flag.Parse()
	outfile := strings.Trim(*fileArg, ".csv")
	outFile := outfile + "_results.csv"
	if _, err := os.Stat(outFile); err != nil {
		if os.IsNotExist(err) {
			csvFile, _ = os.Create(outFile)
			writer = csv.NewWriter(csvFile)
			writer.Write([]string{"user_id", "request_id", "status"})
		}
	} else {
		csvFile, _ = os.OpenFile(outFile, os.O_APPEND|os.O_WRONLY, 0600)
		writer = csv.NewWriter(csvFile)
	}
	byteData, err := ioutil.ReadFile(*fileArg)
	if err != nil {
		panic(err)
	}

	if err := gocsv.UnmarshalBytes(byteData, &users); err != nil {
		panic(err)
	}
	worker := event.NewWorker(writer, users, *startIndex)
	master := rate_limit.RateLimitExecutor{Limit: *rlArg, Concurrency: *concurrency}
	master.Execute(worker)
}
