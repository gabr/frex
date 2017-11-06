package main

import (
	"fmt"
	"os"
	"regexp"
)

const (
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
	regex   *regexp.Regexp
	replace string
	files   []string
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
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return args, fmt.Errorf("File not found: %q", filePath)
		}
	}

	return args, nil
}

func main() {
	_, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
