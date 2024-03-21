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

func main() {
	flag.BoolVar(&countBytes, "c", false, "print the bytes count")
	flag.BoolVar(&countLines, "l", false, "print the lines count")
	flag.Parse()

	var bytesResults []int
	var linesResults []int
	var filenames []string

	for _, filename := range flag.Args() {
		file, err := os.Open(filename)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			break
		}
		filenames = append(filenames, filename)
		defer file.Close()

		bytesCount := 0
		linesCount := 0

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanBytes)

		for scanner.Scan() {
			b := scanner.Bytes()
			if string(b[0]) == "\n" {
				linesCount += 1
			}
			bytesCount += 1
		}

		bytesResults = append(bytesResults, bytesCount)
		linesResults = append(linesResults, linesCount)
	}

	if len(filenames) > 0 {
		var maxValue int
		if countLines {
			maxValue = slices.Max(linesResults)
		} else if countBytes {
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

			if countBytes {
				output += fmt.Sprintf("%*d ", colSize, bytesResults[i])
			}
			fmt.Printf("%s%s\n", output, filename)
		}
	}
}
