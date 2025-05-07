package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	pattern := flag.String("pattern", "", "Search pattern")
	isRegex := flag.Bool("regex", false, "a bool")

	flag.Parse()

	fmt.Println("Grepping for", *pattern, *isRegex)

	if *pattern == "" {
		fmt.Println("Please provide a search pattern using -pattern flag")
		os.Exit(1)
	}

	// Search in provided files
	for _, filename := range flag.Args() {
		searchFile(*pattern, *isRegex, filename)
	}
}

func searchFile(pattern string, isRegex bool, filename string) {
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

		if isRegex {
			// Regex matching
			reg, err := regexp.Compile(pattern)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error compiling regex %s: %v\n", pattern, err)
			}

			if reg.MatchString(line) {
				fmt.Printf("%s:%d:%s\n", filename, lineNum, line)
			}
		} else {
			if strings.Contains(line, pattern) {
				fmt.Printf("%s:%d:%s\n", filename, lineNum, line)
			}
		}
		lineNum++
	}

	err = scanner.Err()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filename, err)
	}
}
