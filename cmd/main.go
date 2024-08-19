package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dirName, oldText, newText := parseArgs()
	fmt.Println(dirName, oldText, newText)
	err := run(dirName, oldText, newText)
	if err != nil {
		fmt.Print(err)
	}

}

func run(path string, newText string, textToReplace string) error {
	absPath, err := filepath.Rel(filepath.Base(""), "./logs")
	if err != nil {
		return err
	}
	logFile, err := os.CreateTemp(absPath, "log.*.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(logFile.Name())
	logger := log.New(logFile, "", log.LstdFlags)
	files, e := os.ReadDir(path)
	if e != nil {
		_, errLog := logFile.WriteString("Something went wrong, check path")
		return errLog
	}
	for _, file := range files {
		err = replacer(file, path, newText, textToReplace, logger)
	}
	return err
}

func replacer(file fs.DirEntry, path string, newText string, textToReplace string, logger *log.Logger) error {
	if !file.IsDir() {
		pathStr := []string{path, file.Name()}
		fmt.Println(file.Name())
		name, errorOpFile := os.OpenFile(strings.Join(pathStr, "/"), os.O_RDWR, 0644)
		if errorOpFile != nil {
			logger.Printf("File reading error, path: %s, file name: %s", path, name)
		}
		scanner := bufio.NewScanner(name)
		var lines []string
		lineNumber := 1
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, textToReplace) {
				var newStr string
				ctr := strings.Count(line, textToReplace)
				for i := ctr; i > 0; i-- {
					index := strings.Index(line, textToReplace)
					newStr = strings.Replace(line, textToReplace, newText, 1)
					logger.Printf("file: %s: line %d, pos: %d \n %s -> %s \n", name.Name(), lineNumber, index, getSnippet(newStr, textToReplace, 5), getSnippet(newStr, newText, 5))
				}
				lines = append(lines, newStr)
			}
			lineNumber++
		}
		err := name.Truncate(0)
		_, err = name.Seek(0, 0)
		writer := bufio.NewWriter(name)
		for _, line := range lines {
			writer.WriteString(line + "\n")
		}
		err = writer.Flush()
		if err != nil {
			log.Fatal(err)
		}
		defer name.Close()
	}

	return nil
}

func getSnippet(oldText string, newText string, snippetLength int) string {
	index := strings.Index(oldText, newText)
	if index == -1 {
		return ""
	}
	start := index - snippetLength
	if start < 0 {
		start = 0
	}
	end := index + len(newText) + snippetLength
	if end > len(oldText) {
		end = len(oldText)
	}
	return oldText[start:end]
}

func parseArgs() (dirName string, oldText string, newText string) {
	flag.StringVar(&dirName, "dirName", "", "path to the directory with files")
	flag.StringVar(&oldText, "replacer", "", "text for replace")
	flag.StringVar(&newText, "searchText", "", "text to be replaced")
	flag.Parse()

	return dirName, oldText, newText
}
