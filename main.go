package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	pattern := flag.String("pattern", "", "Search pattern")
	flag.Parse()

	if *pattern == "" {
		fmt.Println("Please provide a search pattern using -pattern flag")
		os.Exit(1)
	}

	// If no files are provided, read from stdin
	if len(flag.Args()) == 0 {
		searchStdin(*pattern)
		return
	}

	// Search in provided files
	for _, filename := range flag.Args() {
		searchFile(*pattern, filename)
	}
}

func searchFile(pattern, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 1
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, pattern) {
			fmt.Printf("%s:%d:%s\n", filename, lineNum, line)
		}
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filename, err)
	}
}

func searchStdin(pattern string) {
	scanner := bufio.NewScanner(os.Stdin)
	lineNum := 1
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, pattern) {
			fmt.Printf("stdin:%d:%s\n", lineNum, line)
		}
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
	}
}
