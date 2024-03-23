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

type DisplayMode uint32

const (
	LinesMode DisplayMode = 1 << iota
	WordsMode
	CharsMode
	BytesMode
	DefaultMode = LinesMode | WordsMode | BytesMode
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
	maxValue := int64(1)
	for _, fs := range r.FilesStats {
		maxValue = max(maxValue, fs.Bytes, fs.Lines, fs.Words, fs.Chars)
	}
	columnSize := (int)(math.Floor(math.Log10((float64)(maxValue)))) + 1
	return columnSize
}

func (r *Result) Display(mode DisplayMode) string {
	LinesMode := mode&LinesMode > 0
	WordsMode := mode&WordsMode > 0
	BytesMode := mode&BytesMode > 0
	CharsMode := mode&CharsMode > 0

	result := bytes.Buffer{}
	columnSize := r.getColumnSize()

	for _, fs := range r.FilesStats {
		if LinesMode {
			result.WriteString(fmt.Sprintf("%*d ", columnSize, fs.Lines))
		}
		if WordsMode {
			result.WriteString(fmt.Sprintf("%*d ", columnSize, fs.Words))
		}
		if CharsMode {
			result.WriteString(fmt.Sprintf("%*d ", columnSize, fs.Chars))
		}
		if BytesMode {
			result.WriteString(fmt.Sprintf("%*d ", columnSize, fs.Bytes))
		}
		result.WriteString(fmt.Sprintf("%s\n", fs.Filename))
	}
	return result.String()
}

func GetFileStats(name string, f *os.File) FileStats {
	var bytesCnt int64
	var linesCnt int64
	var wordsCnt int64
	var charsCnt int64

	locale := os.Getenv("LC_CTYPE")
	if locale == "" {
		locale = "UTF-8"
	}

	// TODO: support other character encodings
	isMultiBytes := strings.Contains(locale, "UTF")

	r := bufio.NewReader(f)
	for {
		data, err := r.ReadBytes('\n')
		bytesCnt += int64(len(data))
		wordsCnt += int64(len(bytes.Fields(data)))
		if isMultiBytes {
			charsCnt += int64(utf8.RuneCount(data))
		} else {
			charsCnt += 1
		}
		if err != nil {
			break
		}
		linesCnt += 1
	}

	return FileStats{
		Lines:    linesCnt,
		Words:    wordsCnt,
		Chars:    charsCnt,
		Bytes:    bytesCnt,
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

func main() {
	bytesFlag := flag.Bool("c", false, "print the bytes ")
	linesFlag := flag.Bool("l", false, "print the lines ")
	wordsFlag := flag.Bool("w", false, "print the words ")
	charsFlag := flag.Bool("m", false, "print the characters ")
	flag.Parse()

	var mode DisplayMode
	if *linesFlag {
		mode |= LinesMode
	}
	if *wordsFlag {
		mode |= WordsMode
	}
	if *charsFlag {
		mode |= CharsMode
	}
	if *bytesFlag {
		mode |= BytesMode
	}

	// default behavior of wc
	if mode == 0 {
		mode = DefaultMode
	}

	result := Run(flag.Args())
	fmt.Print(result.Display(mode))
}
