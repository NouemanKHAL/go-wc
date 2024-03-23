package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"unicode/utf8"
)

type FileStats struct {
	Bytes    int64
	Lines    int64
	Words    int64
	Chars    int64
	Filename string
}

type Result struct {
	FilesStats []FileStats
}

func (r *Result) getColumnSize() int {
	var maxValue int64
	for _, fs := range r.FilesStats {
		maxValue = max(fs.Bytes, fs.Lines, fs.Words, fs.Chars)
	}
	// to avoid -inf values when maxValue is 0
	maxValue = max(maxValue, 1)
	colSize := (int)(math.Floor(math.Log10((float64)(maxValue)))) + 1
	return colSize
}

func (r *Result) String() string {
	result := bytes.Buffer{}
	columnSize := r.getColumnSize()
	for _, fs := range r.FilesStats {
		if countLines {
			result.WriteString(fmt.Sprintf("%*d ", columnSize, fs.Lines))
		}
		if countWords {
			result.WriteString(fmt.Sprintf("%*d ", columnSize, fs.Words))
		}
		if countChars {
			result.WriteString(fmt.Sprintf("%*d ", columnSize, fs.Chars))
		}
		if countBytes {
			result.WriteString(fmt.Sprintf("%*d ", columnSize, fs.Bytes))
		}
		result.WriteString(fmt.Sprintf("%s\n", fs.Filename))
	}
	return result.String()
}

func GetFileStats(name string, f *os.File) FileStats {
	var bytesCount int64
	var linesCount int64
	var wordsCount int64
	var charsCount int64

	locale := os.Getenv("LC_CTYPE")
	if locale == "" {
		locale = "UTF-8"
	}

	// TODO: support other character encodings
	isMultiBytes := strings.Contains(locale, "UTF")

	r := bufio.NewReader(f)
	for {
		data, err := r.ReadBytes('\n')
		bytesCount += int64(len(data))
		wordsCount += int64(len(bytes.Fields(data)))
		if isMultiBytes {
			charsCount += int64(utf8.RuneCount(data))
		} else {
			charsCount += 1
		}
		if err != nil {
			break
		}
		linesCount += 1
	}

	return FileStats{
		Lines:    linesCount,
		Words:    wordsCount,
		Chars:    charsCount,
		Bytes:    bytesCount,
		Filename: name,
	}
}

func Run(args []string) Result {
	var result Result
	if len(args) == 0 {
		fs := GetFileStats("", os.Stdin)
		result.FilesStats = append(result.FilesStats, fs)
	} else {
		for _, filename := range args {
			file, err := os.Open(filename)
			if err != nil {
				os.Stderr.WriteString(err.Error())
				break
			}
			defer file.Close()

			fs := GetFileStats(filename, file)
			result.FilesStats = append(result.FilesStats, fs)
		}
	}
	return result
}

var countBytes bool
var countLines bool
var countWords bool
var countChars bool

func main() {
	flag.BoolVar(&countBytes, "c", false, "print the bytes count")
	flag.BoolVar(&countLines, "l", false, "print the lines count")
	flag.BoolVar(&countWords, "w", false, "print the words count")
	flag.BoolVar(&countChars, "m", false, "print the characters count")
	flag.Parse()

	// default behavior of wc
	if !(countBytes || countChars || countLines || countWords) {
		countBytes = true
		countLines = true
		countWords = true
	}

	result := Run(flag.Args())
	fmt.Print(result.String())
}
