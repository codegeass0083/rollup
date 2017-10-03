package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	var err error
	// read from stdin
	input, err := ioutil.ReadAll(os.Stdin)
	checkErr(err)
	lines := strings.Split(string(input), "\n")
	inputColumnString := lines[0]
	inputColumns := strings.Split(inputColumnString, "\t")
	var outputFirstLine string
	var outputColumns []string
	// if no arguments are specified, the list of columns is assumed to be all but the final column
	if len(os.Args) == 1 {
		outputColumns = inputColumns[:len(inputColumns)-1]
		outputFirstLine = inputColumnString
	} else {
		outputColumns = os.Args[1:]
		outputFirstLine = strings.Join(outputColumns, "\t") + "\tvalue"
	}
	printOutput(outputFirstLine)
	records := lines[1:]
	// convert records slice to map and do initial aggregation
	recordMap := make(map[string]int64)
	for _, record := range records {
		key, value := generateKV(record, outputColumns, inputColumns)
		if _, ok := recordMap[key]; ok {
			recordMap[key] += value
		} else {
			recordMap[key] = value
		}
	}
	printMap(recordMap, 1)
	// aggregating over all the prefix of output columns and print to stdout
	for i := 0; i < len(outputColumns)-1; i++ {
		recordMap = rollupAndPrint(recordMap, i+2)
	}
	rollupAndPrint(recordMap, len(outputColumns))
}

func rollupAndPrint(recordMap map[string]int64, tabCount int) map[string]int64 {
	ret := make(map[string]int64)
	for k, v := range recordMap {
		key := removeLastColumnFromKey(k)
		if _, ok := ret[key]; ok {
			ret[key] += v
		} else {
			ret[key] = v
		}
	}
	printMap(ret, tabCount)
	return ret
}

func removeLastColumnFromKey(key string) string {
	keySlice := strings.Split(key, "\t")
	return strings.Join(keySlice[:len(keySlice)-1], "\t")
}

func generateKV(record string, output, input []string) (string, int64) {
	keyIndex := generateKeyIndex(output, input)
	recordSlice := strings.Split(record, "\t")
	recordLen := len(recordSlice)
	valueIndex := recordLen - 1
	if recordSlice[valueIndex] == "" {
		// ignore the trailing tab
		valueIndex = recordLen - 2
	}
	value, err := strconv.ParseInt(recordSlice[valueIndex], 10, 64)
	checkErr(err)
	var buffer bytes.Buffer
	for _, i := range keyIndex {
		buffer.WriteString(recordSlice[i] + "\t")
	}
	key := buffer.String()
	return key[:len(key)-1], value
}

func generateKeyIndex(output, input []string) []int {
	ret := make([]int, len(output))
	for i := range ret {
		ret[i] = -1
	}
	for index, outputStr := range output {
		for i, inputStr := range input {
			if inputStr == outputStr {
				ret[index] = i
			}
		}
		if ret[index] == -1 {
			fmt.Printf("Argument %s does not match any column from input file\n", outputStr)
			os.Exit(1)
		}
	}
	return ret
}

func printMap(recordMap map[string]int64, tabCount int) {
	for k, v := range recordMap {
		printOutput(k + getTabs(tabCount) + strconv.FormatInt(v, 10))
	}
}

func getTabs(count int) string {
	var ret bytes.Buffer
	for i := 0; i < count; i++ {
		ret.WriteByte('\t')
	}
	return ret.String()
}

func printOutput(line string) {
	fmt.Println(line)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
