package schema

import (
	"strconv"
	"time"
)

var SubtableUserIdNum int64 = 10000000

// 根据用户id进行分表
func SubtableUserId(userId int64, tableName string) string {
	a := userId / SubtableUserIdNum
	b := userId % SubtableUserIdNum
	if b > 0 {
		a = a + 1
	}
	return tableName + "_" + strconv.FormatInt(a, 10)
}

// 根据时间进行分表（月）根据时间判断出该条数据属于那一月的数据
func SubtableTimeMonth(strTime int64, tableName string) string {
	t := time.Unix(strTime, 0)
	year, month, _ := t.Date()
	return tableName + "_" + strconv.Itoa(year) + "年_" + strconv.Itoa(int(month)) + "月"
}

// 根据时间进行分表（日）根据时间判断出该条数据属于那一天的数据
func SubtableTimeDay(strTime int64, tableName string) string {
	t := time.Unix(strTime, 0)
	year, month, day := t.Date()
	return tableName + "_" + strconv.Itoa(year) + "年_" + strconv.Itoa(int(month)) + "月_" + strconv.Itoa(day) + "日"
}

// 根据时间进行分表（周）根据时间判断出该条数据属于那一周的数据（根据周一00:00为表的起始时间，周天23:59为表的结束时间）
func SubtableTimeWeek(strTime int64, tableName string) string {
	t := time.Unix(strTime, 0)
	weekday := t.Weekday()
	monday := t.AddDate(0, 0, -int(weekday-1))
	year, month, day := monday.Date()
	return tableName + "_" + strconv.Itoa(year) + "年_" + strconv.Itoa(int(month)) + "月_" + strconv.Itoa(day) + "日"
}
