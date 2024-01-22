package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
	"flag"
	"log"
)

type FileStats struct {
	ByteCount      int
	CharacterCount int
	LineCount      int
	WordCount      int
}

func calculateFileStats(fileName string) (FileStats, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return FileStats{}, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var stats FileStats
	var inWord bool

	for {
		r, size, err := reader.ReadRune()
		if err != nil {
			if err != io.EOF {
				return FileStats{}, err
			}
			break
		}

		stats.ByteCount += size
		stats.CharacterCount++

		if r == '\n' {
			stats.LineCount++
		}

		if unicode.IsSpace(r) {
			inWord = false
		} else {
			if !inWord {
				stats.WordCount++
			}
			inWord = true
		}
	}

	if inWord {
		stats.WordCount++
	}

	return stats, nil
}

func main() {
	var printC, printL, printW, printM bool

	// Parse command-line flags
	flag.BoolVar(&printC, "c", false, "Output file size")
	flag.BoolVar(&printL, "l", false, "Output number of lines")
	flag.BoolVar(&printW, "w", false, "Output number of words")
	flag.BoolVar(&printM, "m", false, "Output number of characters")
	flag.Parse()

	printAll := !printC && !printW && !printL && !printM

	if len(os.Args) < 2 {
		log.Fatal("no file specified")
	}
	fileName := os.Args[len(os.Args)-1]
	res := ""
	stats, err := calculateFileStats(fileName)
	if err != nil {
		log.Fatal(err)
	}

	if printC || printAll {
		res += fmt.Sprintf("%8d", stats.ByteCount)
	}

	if printW || printAll {
		res += fmt.Sprintf("%8d", stats.WordCount)
	}

	if printL || printAll {
		res += fmt.Sprintf("%8d", stats.LineCount)
	}

	if printM {
		res += fmt.Sprintf("%8d", stats.CharacterCount)
	}
	res += fmt.Sprintf("%16s", fileName)
	fmt.Print(res)
}
