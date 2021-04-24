package parser

import (
	"bufio"
	"io"
	"strings"
)

const (
	scannerBuffer = 128 * 1024
)

func ParseTextLines(lines []string) ([]Generator, error) {
	spec := strings.Join(lines, " ")
	return ParseString(spec)
}

func ParseReader(reader io.Reader) ([]Generator, error) {
	buffer := make([]byte, scannerBuffer)

	scanner := bufio.NewScanner(reader)
	scanner.Buffer(buffer, scannerBuffer)

	res := []string{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		res = append(res, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ParseTextLines(res)
}
