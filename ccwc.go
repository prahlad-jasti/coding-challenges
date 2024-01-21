package main

import "fmt"
import "os"
import "log"
import "io"

func fileLen(fileName string) (int, error) {
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

func main() {
	if len(os.Args) < 2 {
		log.Fatal("no file specified")
	}
	fileName := os.Args[1]
	size, err := fileLen(fileName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(size)
}
