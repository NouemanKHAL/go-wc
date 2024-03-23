package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

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

func Run(cmdArgs []string) Result {
	var result Result

	// no input ->
	if len(cmdArgs) == 0 {
		fs := GetFileStats("", os.Stdin)
		result.FilesStats = append(result.FilesStats, fs)
		return result
	}

	for _, filename := range cmdArgs {
		file, err := os.Open(filename)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			break
		}
		defer file.Close()

		fs := GetFileStats(filename, file)
		result.FilesStats = append(result.FilesStats, fs)
	}

	return result
}

func main() {
	bytesFlag := flag.Bool("c", false, "print the bytes")
	linesFlag := flag.Bool("l", false, "print the lines")
	wordsFlag := flag.Bool("w", false, "print the words")
	charsFlag := flag.Bool("m", false, "print the characters")

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
