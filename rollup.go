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
		os.Exit(1)
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
	// check first line
	if firstLine != "y\tm\td\tvalue" {
		fmt.Println("Input file has wrong format.")
		return
	}
	fmt.Println(firstLine)
	ymMap := make(map[key]int64)
	yMap := make(map[string]int64)
	for scanner.Scan() {
		record := scanner.Text()
		recordSlice := strings.Split(record, "\t")
		// validate data in each row
		y, err := strconv.ParseInt(recordSlice[0], 10, 64)
		if err != nil || y <= 0 {
			continue
		}
		m, err := strconv.ParseInt(recordSlice[1], 10, 64)
		if err != nil || m < 1 || m > 12 {
			continue
		}
		d, err := strconv.ParseInt(recordSlice[2], 10, 64)
		if err != nil || !isValidDay(m, d) {
			continue
		}
		if len(recordSlice) < 4 {
			//fmt.Printf("[DEBUG] Invalid row %s, skipping ..\n", record)
			continue
		}
		if columnStr == "ymd" || columnStr == "" {
			// print [y m d]
			fmt.Println(record)
			checkErr(err)
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
			// print [y m]
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
		// print [y]
		fmt.Println(k + "\t\t\t" + strconv.FormatInt(v, 10))
		checkErr(err)
		sum += v
	}
	// print []
	fmt.Println("\t\t\t" + strconv.FormatInt(sum, 10))
	checkErr(err)
}

func isValidDay(month, day int64) bool {
	if day < 1 {
		return false
	}
	if month == 2 {
		return day <= 29 // TODO: check leap year
	}
	if month == 4 || month == 6 || month == 9 || month == 11 {
		return day <= 30
	}
	return day <= 31
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
