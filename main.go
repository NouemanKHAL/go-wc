package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"slices"
)

var countBytes bool
var countLines bool
var countWords bool

func main() {
	flag.BoolVar(&countBytes, "c", false, "print the bytes count")
	flag.BoolVar(&countLines, "l", false, "print the lines count")
	flag.BoolVar(&countWords, "w", false, "print the words count")
	flag.Parse()

	// default behavior of wc
	if !countBytes && !countLines && !countWords {
		countBytes = true
		countLines = true
		countWords = true
	}

	var bytesResults []int64
	var linesResults []int64
	var wordsResults []int64
	var filenames []string

	for _, filename := range flag.Args() {
		file, err := os.Open(filename)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			break
		}
		filenames = append(filenames, filename)
		defer file.Close()

		// init counters
		bytesCount := int64(0)
		linesCount := int64(0)
		wordsCount := int64(0)

		if countBytes {
			stat, err := os.Stat(filename)
			if err != nil {
				os.Stderr.WriteString(err.Error())
				break
			}
			bytesCount = stat.Size()
			bytesResults = append(bytesResults, bytesCount)
		}

		if countLines {
			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				linesCount += 1
			}
			linesResults = append(linesResults, linesCount)
		}

		if countWords {
			// rewind the file cursor to the beginning
			file.Seek(0, 0)

			wordScanner := bufio.NewScanner(file)
			wordScanner.Split(bufio.ScanWords)

			for wordScanner.Scan() {
				wordsCount += 1
			}
			wordsResults = append(wordsResults, wordsCount)
		}
	}

	if len(filenames) > 0 {
		var maxValue int64
		if countLines {
			maxValue = slices.Max(linesResults)
		}

		if countWords {
			maxValue = slices.Max(wordsResults)
		}

		if countBytes {
			maxValue = slices.Max(bytesResults)
		}

		// to avoid -inf values when maxValue is 0
		maxValue = max(maxValue, 1)
		colSize := (int)(math.Floor(math.Log10((float64)(maxValue)))) + 1

		for i, filename := range filenames {
			output := ""

			if countLines {
				output += fmt.Sprintf("%*d ", colSize, linesResults[i])
			}
			if countWords {
				output += fmt.Sprintf("%*d ", colSize, wordsResults[i])
			}
			if countBytes {
				output += fmt.Sprintf("%*d ", colSize, bytesResults[i])
			}
			fmt.Printf("%s%s\n", output, filename)
		}
	}
}
