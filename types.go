package main

import (
	"bytes"
	"fmt"
	"math"
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

	// column size is the length of the largest value in the output
	// columnSize = floor(log10(max)) + 1
	columnSize := (int)(math.Floor(math.Log10((float64)(maxValue)))) + 1
	return max(7, columnSize)
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
