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

func main() {
	flag.BoolVar(&countBytes, "c", false, "print the bytes count")
	flag.Parse()

	var bytesResults []int
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

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanBytes)

		for scanner.Scan() {
			bytesCount += 1
		}

		bytesResults = append(bytesResults, bytesCount)
	}

	if len(filenames) > 0 {
		maxValue := slices.Max(bytesResults)
		colSize := (int)(math.Floor(math.Log10((float64)(maxValue)))) + 1

		for i, byteCount := range bytesResults {
			fmt.Printf("%*d %s\n", colSize, byteCount, filenames[i])
		}
	}
}
