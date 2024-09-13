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
	HttpStatusCode HttpStatusCode
	ErrorCode      ErrorCode
	Reason         string
}

// ResponseErr 响应错误
type ResponseErr struct {
	Code   string `json:"code"`
	Reason string `json:"reason"`
}

// Error 实现 error 接口
func (e *Error) Error() string {
	return e.Reason
}

// WithCause 添加错误原因
func (e *Error) WithCause(cause string) error {
	ne := NewError(e.HttpStatusCode, e.ErrorCode, fmt.Sprintf("%s: %s", e.Reason, cause))
	return ne
}

func WithCause(err error, cause string) error {
	switch e := err.(type) {
	case *Error:
		return e.WithCause(cause)
	default:
		ne := fmt.Errorf("%s: %w", cause, e)
		return ne
	}
}

// Data 返回错误数据
func (e *Error) Data() (HttpStatusCode, *ResponseErr) {
	return e.HttpStatusCode, &ResponseErr{
		Code:   string(e.ErrorCode),
		Reason: e.Reason,
	}
}

func NewError(httpStatuCode HttpStatusCode, errorCode ErrorCode, reason string) *Error {
	return &Error{
		HttpStatusCode: httpStatuCode,
		ErrorCode:      errorCode,
		Reason:         reason,
	}
}
