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

	if len(os.Args) < 4 {
		showHelp()
		return
	}

	searchPath := os.Args[1]
	searchExpression := os.Args[2]
	filesExtensions := strings.Split(os.Args[3], ",")

	searchFilesByExtension(searchPath, searchExpression, filesExtensions)
}

// walk by directories looking for files with the extesions suffix
func searchFilesByExtension(searchPath string, searchExpression string, filesExtensions []string) {

	const jobsCapacity = 4

	fileSearch := make(chan string, jobsCapacity*2)
	done := make(chan bool)

	for i := 1; i <= jobsCapacity; i++ {
		go findInFilesConsumer(fileSearch, done, searchExpression)
	}

	oError := filepath.Walk(searchPath,
		func(path string, _ os.FileInfo, oError error) error {

			if oError != nil {
				return oError
			}

			for _, extension := range filesExtensions {
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
func findInFilesConsumer(fileSearch <-chan string, done chan<- bool, searchExpression string) {

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
			if strings.Contains(strings.ToLower(currentLineText), strings.ToLower(searchExpression)) {
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

	for _, oFileLine := range oFoundFile.lines {
		fmt.Println(oFoundFile.pathName+":", oFileLine.line, "->", strings.TrimSpace(oFileLine.content))
	}

}

// show syntax help message to user
func showHelp() {

	fmt.Println("vine: Recursive find-in-files.\n")

	fmt.Println("Syntax: vine <searchPath> <searchExpression> <filesExtensions>")

	fmt.Println("<searchPath>: The starting directory for the search.")
	fmt.Println("<searchExpression>: The string expression to search for.")
	fmt.Println("<filesExtensions>: File extensions to be searched. ")
	fmt.Println("You can use multiple extensions, separating them by ','")
}
