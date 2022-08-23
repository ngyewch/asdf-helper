package util

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

func Scan(r io.Reader, lineHandler func(line string) error) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		idx := strings.Index(line, "#")
		if idx >= 0 {
			line = line[0:idx]
		}
		space := regexp.MustCompile(`\s+`)
		line = space.ReplaceAllString(line, " ")
		line = strings.TrimSpace(line)
		err := lineHandler(line)
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

func ScanFile(filename string, lineHandler func(line string) error) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return Scan(f, lineHandler)
}
