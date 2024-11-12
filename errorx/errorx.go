package errorx

import (
	"fmt"
	"net/http"
)

// HttpStatusCode 表示 HTTP 状态码
type HttpStatusCode int

const (
	// ErrBadRequest 400: 表示客户端发送的请求有错误
	ErrBadRequest HttpStatusCode = http.StatusBadRequest

	// ErrUnauthorized 401: 表示请求需要认证
	ErrUnauthorized HttpStatusCode = http.StatusUnauthorized

	// ErrForbidden 403: 表示服务器拒绝请求
	ErrForbidden HttpStatusCode = http.StatusForbidden

	// ErrNotFound 404: 表示请求的资源不存在
	ErrNotFound HttpStatusCode = http.StatusNotFound

	// ErrMethodNotAllowed 405: 表示请求方法不允许
	ErrMethodNotAllowed HttpStatusCode = http.StatusMethodNotAllowed

	// ErrNotAcceptable 406: 表示请求的资源不支持请求的格式
	ErrNotAcceptable HttpStatusCode = http.StatusNotAcceptable

	// ErrInternalServerError 500: 表示服务器内部错误
	ErrInternalServerError HttpStatusCode = http.StatusInternalServerError
)

// ErrorCode 表示错误码
type ErrorCode string

type Error struct {
	httpStatusCode HttpStatusCode
	errorCode      ErrorCode
	reason         string
	cause          error
}

// ResponseErr 响应错误
type ResponseErr struct {
	Code   string `json:"code"`
	Reason string `json:"reason"`
}

func (e *Error) ErrorCode() string {
	return string(e.errorCode)
}

// Error 实现 error 接口
func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.reason, e.cause)
	}
	return e.reason
}

// Is 判断错误是否相等
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.errorCode == t.errorCode && e.httpStatusCode == t.httpStatusCode
}

func (e *Error) HttpStatusCode() int {
	return int(e.httpStatusCode)
}

// WithCause 添加错误原因
func (e *Error) WithCause(cause string) error {
	return &Error{
		httpStatusCode: e.httpStatusCode,
		errorCode:      e.errorCode,
		reason:         e.reason,
		cause:          fmt.Errorf(cause),
	}
}

func WithCause(err error, cause string) error {
	if err == nil {
		return nil
	}
	switch e := err.(type) {
	case *Error:
		return &Error{
			httpStatusCode: e.httpStatusCode,
			errorCode:      e.errorCode,
			reason:         e.reason,
			cause:          fmt.Errorf(cause),
		}
	default:
		return fmt.Errorf("%s: %w", cause, err)
	}
}

// Data 返回错误数据
func (e *Error) Data() *ResponseErr {
	reason := e.reason
	if e.cause != nil {
		reason = e.Error()
	}
	return &ResponseErr{
		Code:   string(e.errorCode),
		Reason: reason,
	}
}

// NewError 创建一个新的错误
func NewError(httpStatuCode HttpStatusCode, errorCode ErrorCode, reason string) *Error {
	return &Error{
		httpStatusCode: httpStatuCode,
		errorCode:      errorCode,
		reason:         reason,
	}
}
