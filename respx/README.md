# Respx

响应处理库，支持自定义响应码，响应信息，HTTP 状态码, 事件驱动

## 使用

### 事件驱动

```go
func init() {
    respx.RegisterEvent(respx.BeforeResponse, func(w respx.ResponseWriter, data any) {
        // 在响应发送前执行一些操作，比如日志记录
        log.Printf("准备发送响应: %v", data)
    })

    respx.RegisterEvent(respx.AfterResponse, func(w respx.ResponseWriter, data any) {
        // 在响应发送后执行一些操作，比如性能监控
        log.Printf("响应已发送: %v", data)
    })
}
```

### 自定义响应码
