package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
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

	// Create a channel to collect results
	results := make(chan string)
	var wg sync.WaitGroup

	// Search in provided files concurrently
	for _, filename := range flag.Args() {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			searchFileParallel(*pattern, file, results)
		}(filename)
	}

	// Close results channel when all searches are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Print results as they come in
	for result := range results {
		fmt.Print(result)
	}
}

func searchFileParallel(pattern, filename string, results chan<- string) {
	file, err := os.Open(filename)
	if err != nil {
		results <- fmt.Sprintf("Error opening file %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 1
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, pattern) {
			results <- fmt.Sprintf("%s:%d:%s\n", filename, lineNum, line)
		}
		lineNum++
	}

	err = scanner.Err()
	if err != nil {
		results <- fmt.Sprintf("Error reading file %s: %v\n", filename, err)
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

	err := scanner.Err()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
	}
}
