// package main, implementação do vine find in files
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var jobsCapacity int

type fileLine struct {
	line    int
	content string
}

type foundFile struct {
	pathName string
	lines    []fileLine
}

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
				if strings.HasSuffix(strings.ToLower(path), strings.ToLower(extension)) {
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

		var linesFound []fileLine

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
			if strings.Contains(strings.ToLower(currentLineText), strings.ToLower(search)) {
				linesFound = append(linesFound, fileLine{line: line, content: currentLineText})
			}

			line++
		}

		oFile.Close()

		if oError := scanner.Err(); oError != nil {
			log.Println(oError)
		} else {
			if len(linesFound) > 0 {
				showResult(foundFile{pathName: filePath, lines: linesFound})
			}
		}
	}

	done <- true

}

// showing the result to user
func showResult(oFoundFile foundFile) {

	fmt.Println(oFoundFile.pathName)

	for _, oFileLine := range oFoundFile.lines {
		fmt.Println("  ", oFileLine.line, "->", strings.TrimSpace(oFileLine.content))
	}

}
