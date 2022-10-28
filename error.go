package utils

import (
	"context"
	"runtime"

	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

// error 错误值对象
type Error struct {
	Function string       `json:"function"`
	Line     int          `json:"line"`
	Errmsg   string       `json:"errmsg"`
	Param    interface{}  `json:"param"`
	level    logrus.Level `json:"-"`
}

// New
func NewError(level logrus.Level, errmsg string, params ...interface{}) *Error {
	return &Error{
		level:  level,
		Errmsg: errmsg,
		Param:  If(len(params) > 0, params[0], nil),
	}
}

// 错误记录
func (e *Error) String() string {
	strErr, _ := jsoniter.MarshalToString(e)
	return strErr
}

// Error 获取调用者的方法名和调用行
func (e *Error) Error(skip int) string {
	pc, _, line, ok := runtime.Caller(skip)
	if !ok {
		return e.String()
	}
	e.Function = runtime.FuncForPC(pc).Name()
	e.Line = line
	return e.String()
}

// Logger 记录错误日志
func (e *Error) Logger(ctx context.Context) *Error {
	logrus.WithContext(ctx).Logf(e.level, e.Error(2))
	return e
}

// CallerStacks 获取堆栈列表, skip表示堆栈打印到第几层
func CallerStacks(skip int) []*Error {
	var frame runtime.Frame
	pcSlice := make([]uintptr, 50)
	count := runtime.Callers(skip, pcSlice)
	pcSlice = pcSlice[:count]
	frames := runtime.CallersFrames(pcSlice)
	errList := make([]*Error, 0, count)
	for more := count > 0; more; {
		frame, more = frames.Next()
		errList = append(errList, &Error{
			Function: frame.Function,
			Line:     frame.Line,
		})
	}
	return errList
}
