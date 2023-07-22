package output

type LinesData struct {
	StartLine int
	Bytes     []byte
}

type Writer interface {
	LineDataBatchInsert(linesDataBatch ...*LinesData) error
}
