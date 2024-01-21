package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func fileSize(fileName string) (int, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	var total int
	data := make([]byte, 2048)
	for {
		count, err := f.Read(data)
		total += count
		if err != nil {
			if err != io.EOF {
				return 0, err
			}
			break
		}
	}
	return total, nil
}

func fileLines(fileName string) (int, int, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return 0, 0, err
	}
	fileScanner := bufio.NewScanner(f)
	defer f.Close()
	var numLines int
	var numWords int
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		numLines += 1
		line := fileScanner.Text()
		numWords += len(strings.Fields(line))
	}
	return numLines, numWords, nil
}

func countCharacters(fileName string) (int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	charCount := 0
	for {
		_, _, err := reader.ReadRune()
		if err != nil {
			break
		}
		charCount += 1
	}

	return charCount, nil
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
	if printC || printAll {
		size, errSize := fileSize(fileName)
		if errSize != nil {
			log.Fatal(errSize)
		}
		res += fmt.Sprintf("%8d", size)
	}

	if printM {
		chars, errChars := countCharacters(fileName)
		if errChars != nil {
			log.Fatal(errChars)
		}
		res += fmt.Sprintf("%8d", chars)
	}

	if printW || printL || printAll {
		lines, words, errLines := fileLines(fileName)
		if errLines != nil {
			log.Fatal(errLines)
		}
		if printW || printAll {
			res += fmt.Sprintf("%8d", words)
		}
		if printL || printAll {
			res += fmt.Sprintf("%8d", lines)
		}
	}
	res += fmt.Sprintf("%16s", fileName)
	fmt.Print(res)
}
