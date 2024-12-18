package excel

import (
	"io"

	"github.com/xuri/excelize/v2"
)

type ExcelWriter struct {
	f        *excelize.File
	sw       *excelize.StreamWriter
	Filename string
	Headers  []string
}

func NewExcelWriter() (*ExcelWriter, error) {
	f := excelize.NewFile()
	sw, err := f.NewStreamWriter("Sheet1")
	if err != nil {
		return nil, err
	}

	return &ExcelWriter{
		f:  f,
		sw: sw,
	}, nil
}

func (w *ExcelWriter) SetHeaders(headers []string) error {
	w.Headers = headers
	values := make([]interface{}, len(headers))
	for index, item := range headers {
		values[index] = excelize.Cell{Value: item}
	}
	return w.sw.SetRow("A1", values)
}

func (w *ExcelWriter) AddRow(rowIndex int, data []interface{}) error {
	row := make([]interface{}, len(w.Headers))
	for index := 0; index < len(w.Headers); index++ {
		row[index] = data[index]
	}
	cell, err := excelize.CoordinatesToCellName(1, rowIndex)
	if err != nil {
		return err
	}
	return w.sw.SetRow(cell, row)
}

func (w *ExcelWriter) Flush() error {
	return w.sw.Flush()
}

func (w *ExcelWriter) Write(wr io.Writer) error {
	return w.f.Write(wr)
}

func (w *ExcelWriter) Close() error {
	return w.f.Close()
}
