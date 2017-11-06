package main

import (
	"fmt"
	"os"
	"regexp"
)

const (
	fileBufferSize = 1048576 // 1MB
	minArgumentsCount = 3
	help = `Error: %s

frex: File REgeX replacement

Usage:
  frex regex_pattern_to_replace value_to_replace_to file_path1 file_path2...

Example:

> frex .*ReplaceMe Line2 file1.txt

file1.txt before change    >    file1.txt after change
  Line1                    >      Line1
  LineReplaceMe            >      Line2
  Line3                    >      Line3
`
)

type arguments struct {
	regex      *regexp.Regexp
	replace    string
	filesPaths []string
}

// Parses and validates user cli arguments
func parseArgs(userArgs []string) (arguments, error) {
	var err error
	args := arguments{}
	if userArgs == nil {
		return args, fmt.Errorf(help, "Given userArgs are nil")
	}

	if len(userArgs) < 3 {
		return args, fmt.Errorf(help,
			fmt.Sprintf("Not enought arguments.\nGot %d expected at least %d",
				len(userArgs), minArgumentsCount))
	}

	args.regex, err = regexp.Compile(userArgs[0])
	if err != nil {
		return args, err
	}

	args.replace = userArgs[1]

	for _, filePath := range userArgs[2:] {
		// TODO: Remove duplicates
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return args, fmt.Errorf("File not found: %q", filePath)
		}
		append(args.filePaths
	}

	return args, nil
}

func replaceInFile(regex *regexp.Regexp,
                   replace string,
                   filePath string,
                   end chan bool)
{
	defer (func () { end <- true })()	

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening file %q: %s",
			filePath, err)
		return
	}
	defer file.Close()

	// windows new line: 0D 0A
	// linux   new line: 0A
	// mac     new line: 0D or 0A

	buffer := make([]byte, 0, fileBufferSize)
	fileOff := 0
	for {
		n, err := file.ReadAt(buffer, fileOff)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while reading file %q: %s",
				filePath, err)
			return
		}

		fileOff += n

		// TODO: 

		if n < fileBufferSize {
			if n == 0 { break }

			// TODO: Check the rest of the buffer
		}
	}
}

func main() {
	args, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	end := make(chan bool)
	filesCount := len(args.filesPaths)

	for _, filePath := range args.filesPaths {
		go replaceInFile(args.regex, args.replace, filePath, end)
	}

	// wait for all routines to finish
	for i := 0; i < filesCount; i++ {
		<-end
	}
}
