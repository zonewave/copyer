package main

import (
	"bufio"
	"bytes"
	"os"

	"github.com/cockroachdb/errors"
)

type Output struct {
	FileName  string
	StartLine int
}

func NewOutput(fileName string, startLine int) *Output {
	return &Output{FileName: fileName, StartLine: startLine}
}

func (o *Output) Write(bs []byte) (int, error) {
	err := fileInsertLine(o.FileName, o.StartLine, bs)
	return 0, err
}

func fileInsertLine(filePath string, lineNum int, data []byte) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// read lines
	var lines [][]byte
	currentLine := 1
	for scanner.Scan() {
		line := scanner.Bytes()
		lines = append(lines, line)
		if currentLine == lineNum {
			// add new line
			lines = append(lines, data)
		}
		currentLine++
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	// if lineNum is out of range
	if lineNum > currentLine {
		return errors.Errorf("line %d is out of range", lineNum)
	}

	output := bytes.Join(lines, []byte("\n"))
	err = os.WriteFile(filePath, output, 0644)
	if err != nil {
		return errors.Wrapf(err, "write file %s failed", filePath)
	}

	return nil
}
