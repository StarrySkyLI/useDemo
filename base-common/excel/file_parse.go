package excel

import (
	"bytes"
	"encoding/csv"
	"github.com/xuri/excelize/v2"
	"io"
	"net/http"
)

type User struct {
}

// ParseCsv 解析csv文件
func ParseCsv(content *bytes.Buffer, encoding string, separator rune) ([]*User, error) {
	buffer := bytes.NewReader(content.Bytes())
	reader := csv.NewReader(buffer)
	reader.Comma = rune(separator)
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true
	if encoding == "gbk" {
		reader.ReuseRecord = true
		reader.FieldsPerRecord = -1
	}
	users := make([]*User, 0)
	fieldNames, err := reader.Read()
	if err != nil {
		return nil, err
	}
	for _, field := range fieldNames {
		if !CheckField(&User{}, field) {
			return nil, ErrInvalidField
		}
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		user := &User{}
		for i, field := range record {
			err = SetFieldValue(user, fieldNames[i], field)
			if err != nil {
				return nil, err
			}
		}
		users = append(users, user)
	}
	return users, nil
}

// ParseExcel 解析execl文件
func ParseExcel(content *bytes.Buffer) ([]*User, error) {
	xlsx, err := excelize.OpenReader(content)
	if err != nil {
		return nil, err
	}
	sheet := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	users := make([]*User, 0)
	fieldNames := make([]string, 0)
	for _, row := range rows {
		if len(row) == 0 {
			continue
		}
		if len(fieldNames) == 0 {
			for _, cell := range row {
				if !CheckField(&User{}, cell) {
					return nil, ErrInvalidField
				}
				fieldNames = append(fieldNames, cell)
			}
			continue
		}
		user := &User{}
		for i, cell := range row {
			err = SetFieldValue(user, fieldNames[i], cell)
			if err != nil {
				return nil, err
			}
		}
		users = append(users, user)
	}
	return users, nil
}

// DownloadFile 根据提供的地址获取文件内容到一个bytes.Buffer中
func DownloadFile(url string) (*bytes.Buffer, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
