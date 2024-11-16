package xcode

import (
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Details struct {
	Code int `json:"code"`
}

type XCode interface {
	Error() string
	Code() int
	Message() string
	Details() []interface{}
}

type Code struct {
	code int
	msg  string
}

func (c Code) Error() string {
	if len(c.msg) > 0 {
		return c.msg
	}

	return strconv.Itoa(c.code)
}

func (c Code) Code() int {
	return c.code
}

func (c Code) Message() string {
	return c.Error()
}

func (c Code) Details() []interface{} {
	return nil
}

func (c Code) GRPCStatus() *status.Status {
	st := status.New(codes.Code(c.code), c.msg)
	return st
}

func New(code int, msg string) Code {
	return Code{code: code, msg: msg}
}

func NewRpc(code int, msg string) error {
	return New(code, msg).GRPCStatus().Err()
}
