package wtime

import (
	_ "fmt"
	"time"

	"github.com/jinzhu/now"
)

type TmOpenType int32

const (
	// 按天开启
	TmOpenType_TmOpenTypeDay TmOpenType = 0
	// 按周开启
	TmOpenType_TmOpenTypeWeek TmOpenType = 1
	// 按月开启
	TmOpenType_TmOpenTypeMonth TmOpenType = 2
	// 一周多次
	TmOpenType_TmOpenTypeWeek2 TmOpenType = 3
)

type TmEndType int32

const (
	EndNo   TmEndType = -1 // 时间没到
	EndInit TmEndType = 0  // 初始化
	EndOk   TmEndType = 1  // 到了

)

var TmOpenType_name = map[int32]string{
	0: "TmOpenTypeDay",
	1: "TmOpenTypeWeek",
	2: "TmOpenTypeMonth",
}
var TmOpenType_value = map[string]int32{
	"TmOpenTypeDay":   0,
	"TmOpenTypeWeek":  1,
	"TmOpenTypeMonth": 2,
}

var (
	CurTime = &Time_S{1, 0, 8 * 60 * 60} // 东八区
	Format  = "RFC3339"
)

type Time_S struct {
	Id       int   `bson:"_id"`
	Interval int64 `bson:"interval"`
	Offset   int   `bson:"offset"`
}

func GetNow() int64 {
	return time.Now().Unix() + CurTime.Interval
}
func GetMSec() int64 {
	return ((time.Now().UnixNano() + CurTime.Interval*1e6) / 1e6)
}
func GetNano() int64 {
	return (time.Now().UnixNano() + (CurTime.Interval)*1e6)
}
func GetNowTime() time.Time {
	now := time.Now().Add(time.Duration(CurTime.Interval) * time.Second).In(time.FixedZone("CST", CurTime.Offset))
	return now
}
func GetTime(tm int64) time.Time { // 时区转换
	now := time.Unix(tm, 0).In(time.FixedZone("CST", CurTime.Offset))
	return now
}
func SetLocation(hour int) {
	CurTime.Offset = hour * 60 * 60 // 时区
}
func SetTimeInterval(in int64) {
	CurTime.Interval = in - time.Now().Unix()
}

func GetStartTm(cur int64, tm *struct {
	Type  int32 `json:"type"`
	Class int32 `json:"class"`
	Time  int32 `json:"time"`
}) (ret int64) {
	if cur == 0 {
		cur = GetNow()
	}
	cur2 := time.Unix(cur, 0)
	switch tm.Type {
	case int32(TmOpenType_TmOpenTypeDay):
		ret = StartOfDay(cur2).Add(time.Duration(tm.Time) * time.Minute).Unix()
	case int32(TmOpenType_TmOpenTypeWeek), int32(TmOpenType_TmOpenTypeWeek2):
		ret = StartOfWeek(cur2).AddDate(0, 0, int(tm.Class-1)).Add(time.Duration(tm.Time) * time.Minute).Unix()
	case int32(TmOpenType_TmOpenTypeMonth):
		ret = StartOfMonth(cur2).AddDate(0, 0, int(tm.Class)-1).Add(time.Duration(tm.Time) * time.Minute).Unix()
	}
	return
}

func IsSameWeek(t1, t2 time.Time) bool {
	y1, w1 := now.With(t1).BeginningOfWeek().ISOWeek()
	y2, w2 := now.With(t1).BeginningOfWeek().ISOWeek()
	return y1 == y2 && w1 == w2
}

func StartOfDay(in time.Time) time.Time {
	return now.With(in).BeginningOfDay()
}

func IsSameDay(a, b int64) bool {
	return StartOfDay(GetTime(a)) == StartOfDay(GetTime(b))
}

func StartOfWeek(t time.Time) time.Time {
	return now.With(t).BeginningOfWeek()
}
func EndOfWeek(t time.Time) time.Time {
	return now.With(t).EndOfWeek()
}
func StartOfMonth(t time.Time) time.Time {
	return now.With(t).BeginningOfMonth()
}
func EndOfMonth(t time.Time) time.Time {
	return now.With(t).EndOfMonth()
}

/*
*
@brief UTC + 时区
*/
func ParseTimeToUnixTime(in string) uint64 {
	timeT, _ := time.Parse("2006-01-02 15:04:05", in)
	return uint64(timeT.Unix())
}

func GetSinceDayFromTimeStr(in uint64) int64 {
	// timeT, _ := time.Parse("2006-01-02 15:04:05", in)
	timeT := time.Unix(int64(in), 0)
	return int64(GetNowTime().Sub(timeT).Hours()) / 24
}

// 返回当前时间是当年的第几周
func GetWeek() uint32 {
	t := GetNowTime()
	yearDay := t.YearDay()
	yearFirstDay := t.AddDate(0, 0, -yearDay+1)
	firstDayInWeek := int(yearFirstDay.Weekday())

	// 今年第一周有几天
	firstWeekDays := 1
	if firstDayInWeek != 0 {
		firstWeekDays = 7 - firstDayInWeek + 1
	}
	var week int
	if yearDay <= firstWeekDays {
		week = 1
	} else {
		week = (yearDay-firstWeekDays)/7 + 2
	}
	// fmt.Println("%d第%d周", t.Year(), week)
	return uint32(week)
}
