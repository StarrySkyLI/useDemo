package Stdout

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	sysTime "time"

	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
)

type Api_Stdout struct{}

// 日志等级
const (
	C_Kafka_Log_Topic      = "all-server-log-test2"
	C_Kafka_Log_User_Topic = "all-user-log-test2"
	LEVEL_FATA             = 1
	LEVEL_ERROR            = 2
	LEVEL_WARNING          = 3
	LEVEL_INFO             = 4
	LEVEL_DEBUG            = 6
)

// 日志模式
const (
	MODEL_PRO = iota
	MODEL_INFO
	MODEL_DEV
)

// 调用log的服务器名字
var serverName string
var g_ip string
var Mylog = logrus.New()
var address string
var is_open = false

func init() {
	Mylog.SetLevel(logrus.TraceLevel)
	Mylog.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000000",
		ForceColors:     true,
	})
	Mylog.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
	// ConfigLocalFilesystemLogger("./logs", "log.out", sysTime.Hour*24*30, sysTime.Hour*24)
}

func Write(logPath string, logFileName string, maxAge sysTime.Duration, rotationTime sysTime.Duration) {
	baseLogPaht := path.Join(logPath, logFileName)
	fmt.Println(baseLogPaht)

}

func SetLogName(name string) {
	serverName = name
}

func OpenSendLog(name string, open bool, elkAddress string) {
	serverName = name
	is_open = open
	address = elkAddress
	if is_open {
		ConfigESLogger(address, g_ip, serverName)
	}
}

var g_lvl = 2

func callerPrettyfier() (path string) {
	fname := ""
	pc, path, line, ok := runtime.Caller(g_lvl) // 去掉两层，当前函数和日志的接口函数
	if ok {
		if f := runtime.FuncForPC(pc); f != nil {
			fname = f.Name()
		}
	}
	funcName := lastFname(fname)
	path = getFilePath(path)
	return fmt.Sprintf("%s() %s:%d ", funcName, path, line)
}

func ConfigESLogger(esUrl string, esHOst string, index string) error {

	return nil
}

// 设置日志等级
func SetLogLevel(lev logrus.Level) {
	Mylog.SetLevel(lev)
}

func SetPathLvl(lvl int) { g_lvl = lvl }

/*
*设置日志模式
*参数说明:
*@param:mod		模式 MODEL_PRO:只向日志服务器发送日志  MODEL_INFO:向日志服务器发送日志切输出到控制台 MODEL_DEV:只输出到控制台
 */
func SetLogModel(mod int) error {
	if mod <= MODEL_DEV {
		Mylog.SetOutput(os.Stdout)
		log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
	}
	if mod <= MODEL_INFO {
		src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			Mylog.Error(err.Error())
			return err
		}
		writer := bufio.NewWriter(src)
		Mylog.SetOutput(writer)
		log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
	}
	if mod <= MODEL_PRO {
		src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			Mylog.Error(err.Error())
			return err
		}
		writer := bufio.NewWriter(src)
		log.SetOutput(writer)
		Mylog.SetOutput(writer)
	}
	return nil
}

// 安全执行监听函数
func Listen(f interface{}, callback func(interface{}), param string) {
	fname := ""
	pc, path, line, ok := runtime.Caller(1) // 去掉两层，当前函数和日志的接口函数
	if ok {
		if f := runtime.FuncForPC(pc); f != nil {
			fname = f.Name()
		}
	}
	funcName := lastFname(fname)
	path = getFilePath(path)
	timer := sysTime.NewTicker(sysTime.Millisecond * 500)
	success := make(chan bool)
	start := sysTime.Now()
	count := 0
	go func() {
		callback(f)
		close(success)
	}()
	for {
		select {
		case <-success:
			end := sysTime.Since(start).Nanoseconds() / 1000000.00
			timer.Stop()
			if end >= 500 && end < 1000 {
				Mylog.Info(fmt.Sprintf("执行严重超时 %s %s %d (%dms) %s", path, funcName, line, end, param))
			}
			if end >= 1000 && end < 2000 {
				Mylog.Warn(fmt.Sprintf("执行严重超时 %s %s %d (%dms) %s", path, funcName, line, end, param))
			}
			if end >= 2000 {
				Mylog.Error(fmt.Sprintf("执行严重超时 %s %s %d (%dms) %s", path, funcName, line, end, param))
			}
			return
		case <-timer.C:
			count++
			end := sysTime.Since(start).Nanoseconds() / 1000000.00
			if count >= 10 {
				Mylog.Error(fmt.Sprintf("执行严重超时%d次提醒 %s %s %d (%dms) %s", count, path, funcName, line, end, param))
			} else {
				Mylog.Info(fmt.Sprintf("执行严重超时%d次提醒 %s %s %d (%dms) %s", count, path, funcName, line, end, param))
			}
		}
	}
}

// 计算函数所用时间
func (one *Api_Stdout) TraceParam(param string) func() {
	fname := ""
	pc, path, line, ok := runtime.Caller(1) // 去掉两层，当前函数和日志的接口函数
	if ok {
		if f := runtime.FuncForPC(pc); f != nil {
			fname = f.Name()
		}
	}
	funcName := lastFname(fname)
	path = getFilePath(path)

	start := sysTime.Now()
	return func() {
		end := sysTime.Since(start).Nanoseconds() / 1000000.00
		if end >= 100 && end < 1000 {
			one.Debug("执行严重超时提醒 %s %s %d (%dms) %s", path, funcName, line, end, param)
		}
		if end >= 1000 && end < 2000 {
			one.Warn("执行严重超时提醒 %s %s %d (%dms) %s", path, funcName, line, end, param)
		}
		if end >= 2000 {
			one.Error(nil, "执行严重超时提醒 %s %s %d (%dms) %s", path, funcName, line, end, param)
		}
	}
}

// 计算函数所用时间
func (one *Api_Stdout) Trace() func() {
	fname := ""
	pc, path, line, ok := runtime.Caller(1) // 去掉两层，当前函数和日志的接口函数
	if ok {
		if f := runtime.FuncForPC(pc); f != nil {
			fname = f.Name()
		}
	}
	funcName := lastFname(fname)
	path = getFilePath(path)

	start := sysTime.Now()
	return func() {
		end := sysTime.Since(start).Nanoseconds() / 1000000.00
		if end >= 100 && end < 1000 {
			one.Info("执行严重超时提醒 %s %s %d (%dms)", path, funcName, line, end)
		}
		if end >= 1000 && end < 2000 {
			one.Warn("执行严重超时提醒 %s %s %d (%dms)", path, funcName, line, end)
		}
		if end >= 2000 {
			one.Error(nil, "执行严重超时提醒 %s %s %d (%dms)", path, funcName, line, end)
		}
	}
}

// 计算函数所用时间
func (one *Api_Stdout) TraceInfo(str string) func() {
	fname := ""
	pc, path, line, ok := runtime.Caller(1) // 去掉两层，当前函数和日志的接口函数
	if ok {
		if f := runtime.FuncForPC(pc); f != nil {
			fname = f.Name()
		}
	}
	funcName := lastFname(fname)
	path = getFilePath(path)

	start := sysTime.Now()
	return func() {
		end := sysTime.Since(start).Nanoseconds() / 1000000.00
		if end >= 100 && end < 1000 {
			one.Info("%s 执行严重超时提醒 %s %s %d (%dms)", str, path, funcName, line, end)
		}
		if end >= 1000 && end < 2000 {
			one.Warn("%s 执行严重超时提醒 %s %s %d (%dms)", str, path, funcName, line, end)
		}
		if end >= 2000 {
			one.Error(nil, "%s 执行严重超时提醒 %s %s %d (%dms)", str, path, funcName, line, end)
		}
	}
}

func GetFunctionName(i interface{}, seps ...rune) string {
	// 获取函数名称
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()

	// 用 seps 进行分割
	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})

	// fmt.Println(fields)

	if size := len(fields); size > 0 {
		return fields[size-1]
	}
	return ""
}

func UserInfoLog(userId int64, format string, args ...interface{}) {
	// defer Trace()()
	Mylog.WithFields(logrus.Fields{
		"UserId": userId,
	}).Warn(callerPrettyfier() + fmt.Sprintf(format, args...))
}

// 错误
func (*Api_Stdout) Error(ctx context.Context, args ...interface{}) {
	// defer Trace()()
	args = append(args, string(debug.Stack()))
	Mylog.Error(callerPrettyfier() + fmt.Sprintf(args[0].(string), args[1:]...))
}

// 警告
func (*Api_Stdout) Warning(args ...interface{}) {
	// defer Trace()()
	Mylog.Warn(callerPrettyfier() + fmt.Sprintf(args[0].(string), args[1:]...))
}

func (*Api_Stdout) Fatal(args ...interface{}) {
	// defer Trace()()
	Mylog.Warn(callerPrettyfier() + fmt.Sprintf(args[0].(string), args[1:]...))
}

// 错误
func (*Api_Stdout) Fatalf(format string, args ...interface{}) {
	// defer Trace()()
	Mylog.Error(callerPrettyfier() + fmt.Sprintf(format, args...))
}

// 警告
func (*Api_Stdout) Warn(args ...interface{}) {
	// defer Trace()()
	Mylog.Warn(callerPrettyfier() + fmt.Sprintf(args[0].(string), args[1:]...))
}

// 提示
func (*Api_Stdout) Info(args ...interface{}) {
	// defer Trace()()
	if len(args) < 2 {
		Mylog.Info(callerPrettyfier(), args)
	} else {
		Mylog.Info(callerPrettyfier() + fmt.Sprintf(args[0].(string), args[1:]...))
	}
}

func (*Api_Stdout) Notice(args ...interface{}) {
	// defer Trace()()
	Mylog.Info(callerPrettyfier() + fmt.Sprintf(args[0].(string), args[1:]...))
}
func (*Api_Stdout) Panic(args ...interface{}) {
	// defer Trace()()
	Mylog.Error(callerPrettyfier() + fmt.Sprintf(args[0].(string), args[1:]...))
	panic(args)
}
func (*Api_Stdout) Critical(args ...interface{}) {
	// defer Trace()()
	Mylog.Info(callerPrettyfier() + fmt.Sprintf(args[0].(string), args[1:]...))
}

// 调试
func (*Api_Stdout) Debug(args ...interface{}) {
	// defer Trace()()
	Mylog.Debug(callerPrettyfier() + fmt.Sprintf(args[0].(string), args[1:]...))
}

// 错误
func (*Api_Stdout) Errorf(ctx context.Context, format string, args ...interface{}) {
	// defer Trace()()
	Mylog.Error(callerPrettyfier() + fmt.Sprintf(format, args...))
}

// 警告
func (*Api_Stdout) Warningf(format string, args ...interface{}) {
	// defer Trace()()
	Mylog.Warn(callerPrettyfier() + fmt.Sprintf(format, args...))
}

// 警告
func (*Api_Stdout) Panicf(format string, args ...interface{}) {
	// defer Trace()()
	Mylog.Warn(callerPrettyfier() + fmt.Sprintf(format, args...))
}

// 提示
func (*Api_Stdout) Infof(format string, args ...interface{}) {
	// defer Trace()()
	Mylog.Info(callerPrettyfier() + fmt.Sprintf(format, args...))
}

// 调试
func (*Api_Stdout) Debugf(format string, args ...interface{}) {
	// defer Trace()()
	Mylog.Debug(callerPrettyfier() + fmt.Sprintf(format, args...))
}

func Test(format string, args ...interface{}) {
	// defer Trace()()
	Mylog.Trace(callerPrettyfier() + fmt.Sprintf(format, args...))
}

// func logFormat(msg string) string {
// 	fname := ""
// 	pc, path, line, ok := runtime.Caller(2) // 去掉两层，当前函数和日志的接口函数
// 	if ok {
// 		if f := runtime.FuncForPC(pc); f != nil {
// 			fname = f.Name()
// 		}
// 	}
// 	funcName := lastFname(fname)
// 	path = getFilePath(path)
// 	format := fmt.Sprintf(" %s %s %d ▶ %s", path, funcName, line, msg)
// 	//fmt.Println(format)
// 	return format
// }

func lastFname(fname string) string {
	flen := len(fname)
	n := strings.LastIndex(fname, ".")
	if n+1 < flen {
		return fname[n+1:]
	}
	return fname
}

func getFilePath(path string) string {
	// s := strings.Split(path, "src")
	return path
}
