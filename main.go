package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"math"
	"os"
	"unicode/utf8"
)

var countBytes bool
var countLines bool
var countWords bool
var countChars bool

type Row struct {
	Bytes    int
	Lines    int
	Words    int
	Chars    int
	Filename string
}

type Result struct {
	Rows []Row
}

func (r *Result) getColumnSize() int {
	var maxValue int
	for _, row := range r.Rows {
		maxValue = max(row.Bytes, row.Lines, row.Words, row.Chars)
	}
	// to avoid -inf values when maxValue is 0
	maxValue = max(maxValue, 1)
	colSize := (int)(math.Floor(math.Log10((float64)(maxValue)))) + 1
	return colSize
}

func (r *Result) String() string {
	result := bytes.Buffer{}
	colSize := r.getColumnSize()
	for _, row := range r.Rows {
		output := ""
		if countLines {
			output += fmt.Sprintf("%*d ", colSize, row.Lines)
		}
		if countWords {
			output += fmt.Sprintf("%*d ", colSize, row.Words)
		}
		if countChars {
			output += fmt.Sprintf("%*d ", colSize, row.Chars)
		}
		if countBytes {
			output += fmt.Sprintf("%*d ", colSize, row.Bytes)
		}
		result.WriteString(fmt.Sprintf("%s%s\n", output, row.Filename))
	}
	return result.String()
}

func GetFileStats(name string, f *os.File) Row {
	bytesCount := 0
	linesCount := 0
	wordsCount := 0
	charsCount := 0

	locale := os.Getenv("LC_CTYPE")
	if locale == "" {
		locale = "UTF-8"
	}

	r := bufio.NewReader(f)
	for {
		data, err := r.ReadBytes('\n')
		bytesCount += len(data)
		wordsCount += len(bytes.Fields(data))
		if locale == "UTF-8" {
			charsCount += utf8.RuneCount(data)
		} else {
			charsCount += 1
		}
		if err != nil {
			break
		}
		linesCount += 1
	}

	return Row{
		Lines:    linesCount,
		Words:    wordsCount,
		Chars:    charsCount,
		Bytes:    bytesCount,
		Filename: name,
	}
}

func main() {
	flag.BoolVar(&countBytes, "c", false, "print the bytes count")
	flag.BoolVar(&countLines, "l", false, "print the lines count")
	flag.BoolVar(&countWords, "w", false, "print the words count")
	flag.BoolVar(&countChars, "m", false, "print the characters count")
	flag.Parse()

	// default behavior of wc
	if !countBytes && !countLines && !countWords && !countChars {
		countBytes = true
		countLines = true
		countWords = true
	}

	var result Result
	if len(flag.Args()) == 0 {
		row := GetFileStats("", os.Stdin)
		result.Rows = append(result.Rows, row)
	}

	for _, filename := range flag.Args() {
		file, err := os.Open(filename)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			break
		}
		defer file.Close()

		row := GetFileStats(filename, file)
		result.Rows = append(result.Rows, row)
	}
	fmt.Print(result.String())
}
