# errorx 包

errorx 包提供了一个增强的错误处理机制，特别适用于 HTTP 服务。它定义了自定义的错误类型，包含 HTTP 状态码、错误码和原因，并提供了创建、包装和处理这些错误的方法。

## 类型

### HttpStatusCode

`HttpStatusCode` 是一个表示 HTTP 状态码的类型别名：

type HttpStatusCode int

预定义的 HTTP 状态码常量包括：

- `ErrBadRequest` (400): 表示客户端发送的请求有错误
- `ErrUnauthorized` (401): 表示请求需要认证
- `ErrForbidden` (403): 表示服务器拒绝请求
- `ErrNotFound` (404): 表示请求的资源不存在
- `ErrMethodNotAllowed` (405): 表示请求方法不允许
- `ErrNotAcceptable` (406): 表示请求的资源不支持请求的格式
- `ErrInternalServerError` (500): 表示服务器内部错误

### ErrorCode

`ErrorCode` 是一个表示错误码的类型别名：

type ErrorCode string

### Error

`Error` 结构体表示一个增强的错误，包含 HTTP 状态码、错误码和原因：

type Error struct {
    HttpStatusCode HttpStatusCode
    ErrorCode      ErrorCode
    Reason         string
}

### ResponseErr

`ResponseErr` 结构体用于 JSON 响应中的错误信息：

type ResponseErr struct {
    Code   string `json:"code"`
    Reason string `json:"reason"`
}

## 方法

### Error.Error()

实现了 `error` 接口，返回错误的原因。

### Error.WithCause(cause string) error

为现有的 `Error` 添加额外的错误原因，返回一个新的错误。

### WithCause(err error, cause string) error

为任意错误添加额外的错误原因。如果输入是 `*Error` 类型，它会调用相应的方法；否则，它会创建一个新的错误，将原始错误包装在其中。

### Error.Data() (HttpStatusCode, *ResponseErr)

返回错误的 HTTP 状态码和用于响应的错误数据。

### NewError(httpStatuCode HttpStatusCode, errorCode ErrorCode, reason string) *Error

创建并返回一个新的 `Error` 实例。

## 使用示例
```golang
// 创建一个新的错误
err := errorx.NewError(errorx.ErrBadRequest, "ERR001", "Invalid input")

// 添加额外的错误原因
newErr := errorx.WithCause(err, "Missing required field")

// 获取错误数据用于 HTTP 响应
statusCode, respErr := newErr.Data()

// 使用错误数据
fmt.Printf("Status Code: %d\n", statusCode)
fmt.Printf("Error Code: %s\n", respErr.Code)
fmt.Printf("Error Reason: %s\n", respErr.Reason)
```

这个包设计用于简化 HTTP 服务中的错误处理，提供了一种统一的方式来创建、包装和处理错误，同时保持与 HTTP 状态码的一致性。