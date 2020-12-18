package errcode

import (
	"fmt"
)

// StandardError struct
type StandardError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (err StandardError) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", err.Code, err.Msg)
}

// WithMsg is method to set error msg
func (err StandardError) WithMsg(msg string) StandardError {
	err.Msg = err.Msg + ": " + msg
	return err
}

// New StarandError
func New(code int, msg string) StandardError {
	return StandardError{code, msg}
}

var (
	Success            = StandardError{0, "success"}
	ErrUnknown         = StandardError{-1, "unknown error"}
	ErrReqForbidden    = StandardError{11001, "登录态失效，请重新登录"}
	ErrParam           = StandardError{11003, "接口参数错误"}
	ErrMethodIncorrect = StandardError{11004, "method incorrect"}
	ErrTimeout         = StandardError{11005, "timeout"}
	ErrDbQuery         = StandardError{11006, "数据库操作错误，请重试"}
	ErrDbConnect       = StandardError{11007, "数据库连接错误，请重试"}
	ErrUserPwd         = StandardError{11008, "用户名或密码错误"}
	ErrInternal        = StandardError{11009, "内部错误"}
	ErrRedisOper       = StandardError{11007, "Redis操作失败，请重试"}
)
