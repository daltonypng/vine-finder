// package main, implementação do vine find in files
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var jobsCapacity int

func main() {

	if len(os.Args) < 3 {
		log.Println("This program requires 2 arguments")
		return
	}

	search := os.Args[1]
	extensions := strings.Split(os.Args[2], ",")
	jobsCapacity = 3

	searchFilesByExtension(search, extensions)
}

// walk by directories looking for files to push into consumers
func searchFilesByExtension(search string, extensions []string) {

	fileSearch := make(chan string, jobsCapacity*2)
	done := make(chan bool)

	for i := 1; i <= jobsCapacity; i++ {
		go findInFilesConsumer(fileSearch, done, search)
	}

	oError := filepath.Walk(".",
		func(path string, _ os.FileInfo, oError error) error {

			if oError != nil {
				return oError
			}

			for _, extension := range extensions {
				if strings.HasSuffix(path, extension) {
					fileSearch <- path
				}
			}

			return nil
		})

	if oError != nil {
		log.Println(oError)
	}

	close(fileSearch)

	// waiting for all consumers to get done
	for i := 1; i <= jobsCapacity; i++ {
		<-done
	}
}

// consume all files pushed, to execute the searching
func findInFilesConsumer(fileSearch <-chan string, done chan<- bool, search string) {

	for filePath := range fileSearch {

		oFile, oError := os.Open(filePath)

		if oError != nil {
			log.Println(oError)
			return
		}

		// Splits on newlines by default.
		scanner := bufio.NewScanner(oFile)

		line := 1

		for scanner.Scan() {
			currentLineText := scanner.Text()
			if strings.Contains(currentLineText, search) {
				fmt.Println(filePath + ":" + strconv.Itoa(line) + " -> " + strings.Trim(currentLineText, " "))
			}

			line++
		}

		oFile.Close()

		if oError := scanner.Err(); oError != nil {
			log.Println(oError)
		}
	}

	done <- true

}
