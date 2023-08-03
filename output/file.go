package output

import (
	"bufio"
	"bytes"
	"os"
	"sort"

	"github.com/cockroachdb/errors"
	"github.com/duke-git/lancet/v2/iterator"
)

type File struct {
	FileName string
}

func NewFile(fileName string) *File {
	return &File{FileName: fileName}
}

func (o *File) LineDataBatchInsert(linesDataBatch ...*LinesData) error {
	file, err := os.Open(o.FileName)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	sort.Slice(linesDataBatch, func(i, j int) bool {
		return linesDataBatch[i].StartLine < linesDataBatch[j].StartLine
	})
	newLinesDataItr := iterator.FromSlice(linesDataBatch)
	newLinesData, _ := newLinesDataItr.Next()
	// read lines
	var lines [][]byte
	currentLine := 1
	for scanner.Scan() {
		lines = append(lines, scanner.Bytes())
		for newLinesData != nil && currentLine == newLinesData.StartLine {
			lines = append(lines, newLinesData.Bytes)
			newLinesData, _ = newLinesDataItr.Next()
		}
		currentLine++
	}

	if err = scanner.Err(); err != nil {
		return err
	}
	for newLinesData != nil {
		lines = append(lines, newLinesData.Bytes)
		newLinesData, _ = newLinesDataItr.Next()
	}

	// if newLineNum is out of range
	if newLinesData != nil {
		return errors.Errorf("lineNum %d out of range", newLinesData.StartLine)
	}

	output := bytes.Join(lines, []byte("\n"))
	output = append(output, []byte("\n")...)
	err = os.WriteFile(o.FileName, output, 0644)
	if err != nil {
		return errors.Wrapf(err, "write file %s failed", o.FileName)
	}

	return nil
}
