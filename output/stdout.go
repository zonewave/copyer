package output

import "os"

type Stdout struct {
	*os.File
}

func (s *Stdout) LineDataBatchInsert(linesDataBatch ...*LinesData) error {
	for _, linesData := range linesDataBatch {
		_, err := s.Write(linesData.Bytes)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewStdout() *Stdout {
	return &Stdout{os.Stdout}
}
