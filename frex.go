package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TODO: Base order on the operating system
var newLineSequences []string = []string{
	"\x0D\x0A", // windows
	"\x0A",     // linux
	"\x0D",     // mac
}

const (
	fileBufferSize    = 255
	minArgumentsCount = 3
	help              = `Error: %s

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

	addedFilePaths := make(map[string]bool)
	for _, filePath := range userArgs[2:] {
		// ignore duplicates
		_, added := addedFilePaths[filePath]
		if added {
			fmt.Fprintf(os.Stderr, "Ignoring duplicated path: '%s'\n", filePath)
			continue
		}

		// TODO: Accept wildchards
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return args, fmt.Errorf("File not found: %q", filePath)
		}

		args.filesPaths = append(args.filesPaths, filePath)
		addedFilePaths[filePath] = true
	}

	return args, nil
}

// Finds out closest new line pos of windows/linux/mac new lines sequences.
// Returns both the position and the new line sequence.
//
// Errors:
//   If new line is not found returns -1
//   If the given offset if greater then given string length it will return -2
//   If the given offset is negative then it will return -3
func findOutNewLinePos(s string, offset int) (pos int, seq string) {
	if offset < 0 {
		return -3, ""
	}

	if offset > len(s) {
		return -2, ""
	}

	s = s[offset:]

	minPos := -1
	minPosSeq := ""
	for _, seq := range newLineSequences {
		pos := strings.Index(s, seq)
		if pos != -1 {
			if minPos == -1 {
				minPos = pos
				minPosSeq = seq
			} else if pos < minPos {
				minPos = pos
				minPosSeq = seq
			}
		}
	}

	return minPos, minPosSeq
}

func replaceInFile(regex *regexp.Regexp, replace string, path string, end chan bool) {
	defer (func() { end <- true })()

	fmt.Printf("DEBUG (%q): Started\n", path)
	defer (func() { fmt.Printf("DEBUG (%q): Ended\n", path) })()

	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening file %q: %s", path, err)
		return
	}
	defer file.Close()

	buffer := make([]byte, fileBufferSize, fileBufferSize)
	var fileOff int64 = 0
	for {
		fmt.Printf("DEBUG (%q): In loop\n", path)

		n, err := file.ReadAt(buffer, fileOff)
		fmt.Printf("DEBUG (%q): n: %d/%d Buffer: %v\n", path, n,
			fileBufferSize, buffer)

		// the n == 0 is to not catch EOF until readed all bytes till the end
		if err != nil && n == 0 {
			fmt.Fprintf(os.Stderr, "Error while reading file '%s': %s\n",
				path, err)
			return
		}

		fileOff += int64(n) // TODO: Only if no pattern was found in line

		// TODO: Plan:
		// 1. Read buffer
		// 2. Try to find new line sequence
		//    If found ambiguity write error and end
		// 3. In each line slice check regex and replace if required
		// 4. If replace is required write only this change to file
		//    and continue reading buffer just after the last rune of
		//    the change.
		// 5. Repeat until read less then buffer

		if n < fileBufferSize {
			fmt.Printf("DEBUG (%q): In condition %d < %d\n", path, n, fileBufferSize)
			if n == 0 {
				break
			}

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
