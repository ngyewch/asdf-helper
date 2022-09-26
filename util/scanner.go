package util

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

func Scan(r io.Reader, lineHandler func(line string, comment string) error) error {
	spaceRegex := regexp.MustCompile(`\s+`)
	commentRegex := regexp.MustCompile(`#+`)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		comment := ""
		commentStartLocation := commentRegex.FindStringIndex(line)
		if commentStartLocation != nil {
			orgLine := line
			comment = orgLine[commentStartLocation[0]:]
			line = orgLine[0:commentStartLocation[0]]
		}
		line = spaceRegex.ReplaceAllString(line, " ")
		line = strings.TrimSpace(line)
		//comment = commentRegex.ReplaceAllString(comment, "")
		comment = strings.TrimSpace(comment)
		err := lineHandler(line, comment)
		if err != nil {
			return err
		}
	}
	err := scanner.Err()
	if err != nil {
		return err
	}
	return nil
}

func ScanFile(filename string, lineHandler func(line string, comment string) error) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return Scan(f, lineHandler)
}
