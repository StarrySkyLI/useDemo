package excel

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	TABLE_SPLIT_SIZE = 10000000 //分表大小
	EXPORT_MAX_SIZE  = 10000
)

var ErrInvalidField = errors.New("invalid field")

func StrToTime(str string) (time.Time, error) {
	str = strings.Replace(str, "/", "-", -1)
	reg, _ := regexp.Compile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`)
	if reg.MatchString(str) {
		return time.Parse("2006-01-02 15:04:05", str)
	}
	reg, _ = regexp.Compile(`^\d{4}-\d{2}-\d{2}$`)
	if reg.MatchString(str) {
		return time.Parse("2006-01-02", str)
	}
	return time.Time{}, errors.New("InvalidTimeFormat")
}

func GetTimeRangeOfDay(day string) (time.Time, time.Time, error) {
	t, err := StrToTime(day)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	startTimeStr := t.Format("2006-01-02 00:000:00")
	endTimeStr := t.Format("2006-01-02 23:59:59")
	startTime, _ := StrToTime(startTimeStr)
	endTime, _ := StrToTime(endTimeStr)
	return startTime, endTime, nil
}

func GetSplitTableName(tablePrefix string, id int64) string {
	mod := id % TABLE_SPLIT_SIZE
	num := id / TABLE_SPLIT_SIZE
	if mod > 0 {
		num += 1
	}

	table := ""
	switch tablePrefix {
	default:
		table = tablePrefix + fmt.Sprintf("%d", num)
	}
	return table
}

// 检查某个字段是否存在于结构体中
func CheckField(obj interface{}, field string) bool {
	_, ok := obj.(map[string]interface{})[field]
	return ok
}

// 提供字段名，字段值，为结构体设置字段值
func SetFieldValue(obj interface{}, field string, value interface{}) error {
	if !CheckField(obj, field) {
		return ErrInvalidField
	}
	obj.(map[string]interface{})[field] = value
	return nil
}

// 时间戳处理
func TimeDeal(timestamp int64) int64 {
	// 如果是毫秒转换成秒
	if timestamp > 1000000000000 {
		return timestamp / 1000
	}
	return timestamp
}
