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
	Success              = StandardError{0, "success"}
	ErrUnknown           = StandardError{-1, "unknown error"}
	ErrReqForbidden      = StandardError{11001, "登录态失效，请重新登录"}
	ErrParam             = StandardError{11003, "接口参数错误"}
	ErrMethodIncorrect   = StandardError{11004, "method incorrect"}
	ErrTimeout           = StandardError{11005, "timeout"}
	ErrDbQuery           = StandardError{11006, "数据库操作错误，请重试"}
	ErrDbConnect         = StandardError{11007, "数据库连接错误，请重试"}
	ErrUserPwd           = StandardError{11008, "用户名或密码错误"}
	ErrInternal          = StandardError{11009, "内部错误"}
	ErrRedisOper         = StandardError{11010, "Redis操作失败，请重试"}
	ErrEmail             = StandardError{11011, "邮箱有误，请检查"}
	ErrEmailAlReadyValid = StandardError{11012, "邮箱已验证，不需要重复验证"}
	ErrEmailValidLimit   = StandardError{11013, "邮箱验证次数超过限制"}
	ErrCodeInValid       = StandardError{11014, "code 无效"}
	ErrUserExisted       = StandardError{11015, "用户已经存在"}
	ErrEmailExisted       = StandardError{11016, "邮箱已经被注册"}
)
