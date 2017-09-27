package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type key struct {
	year  string
	month string
}

func main() {
	// args: inputPath, y, m, d
	commandLineArgs := os.Args[1:]
	if len(commandLineArgs) < 1 || len(commandLineArgs) > 4 {
		fmt.Println("Wrong number of arguments. Usage: ./rollup [input_file_path y m d]")
		return
	}
	// read data from input file
	tsvFile, err := os.Open(commandLineArgs[0])
	checkErr(err)
	defer tsvFile.Close()

	// validate column arguments
	columnStr := strings.Join(commandLineArgs[1:], "")
	if columnStr != "y" && columnStr != "ym" && columnStr != "ymd" && columnStr != "" {
		fmt.Println("Invalid input for columns, only accept [y, m, d], [y, m], [y], []")
		os.Exit(1)
	}

	// data processing
	scanner := bufio.NewScanner(tsvFile)
	if !scanner.Scan() {
		fmt.Println("Input is empty.")
		return
	}
	firstLine := scanner.Text()
	// copy first line to stdout
	fmt.Println(firstLine)
	checkErr(err)
	ymMap := make(map[key]int64)
	yMap := make(map[string]int64)
	for scanner.Scan() {
		record := scanner.Text()
		if columnStr == "ymd" || columnStr == "" {
			fmt.Println(record)
			checkErr(err)
		}
		recordSlice := strings.Split(record, "\t")
		if len(recordSlice) < 4 {
			fmt.Printf("Invalid row %s, skipping ..\n", record)
			continue
		}
		recordKey := key{year: recordSlice[0], month: recordSlice[1]}
		val, err := strconv.ParseInt(recordSlice[3], 10, 64)
		//fmt.Printf("[DEBUG] %v\n", val)
		checkErr(err)
		if _, ok := ymMap[recordKey]; ok {
			ymMap[recordKey] += val
		} else {
			ymMap[recordKey] = val
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	for k, v := range ymMap {
		//fmt.Printf("[DEBUG] ymMap k %s %s v %d\n", k.year, k.month, v)
		if columnStr != "y" {
			fmt.Println(k.year + "\t" + k.month + "\t\t" + strconv.FormatInt(v, 10))
			checkErr(err)
		}
		if _, ok := yMap[k.year]; ok {
			yMap[k.year] += v
		} else {
			yMap[k.year] = v
		}
	}
	var sum int64
	for k, v := range yMap {
		fmt.Println(k + "\t\t\t" + strconv.FormatInt(v, 10))
		checkErr(err)
		sum += v
	}
	fmt.Println("\t\t\t" + strconv.FormatInt(sum, 10))
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
