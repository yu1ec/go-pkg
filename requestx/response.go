package requestx

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/launchdarkly/eventsource"
	"github.com/tidwall/gjson"
)

type Response struct {
	resp   *http.Response
	req    *http.Request
	body   []byte
	stream chan []byte
	err    error
}

type ResponseBody []byte

// String 获取body体字符串
func (r ResponseBody) String() string {
	return string(r)
}

// Read 读取指定长度的body体
func (r ResponseBody) Read(length int) []byte {
	if length > len(r) {
		length = len(r)
	}
	return r[:length]
}

// GetContents 获取body体字符串
func (r ResponseBody) GetContents() string {
	return string(r)
}

// GetRequest 获取请求
func (r *Response) GetRequest() *http.Request {
	return r.req
}

// GetBody 获取body体
func (r *Response) GetBody() (ResponseBody, error) {
	return ResponseBody(r.body), r.err
}

// GetParsedBody 获取json格式的body体 gjson.Result
func (r *Response) GetParsedBody() (*gjson.Result, error) {
	pb := gjson.ParseBytes(r.body)
	return &pb, nil
}

func (r *Response) GetStatusCode() int {
	return r.resp.StatusCode
}

// GetReasonPhrase 获取响应状态说明
func (r *Response) GetReasonPhrase() string {
	status := r.resp.Status
	arr := strings.Split(status, " ")
	return arr[1]
}

func (r *Response) IsTimeout() bool {
	if r.err == nil {
		return false
	}
	netErr, ok := r.err.(net.Error)
	if !ok {
		return false
	}

	if netErr.Timeout() {
		return true
	}

	return false
}

func (r *Response) GetHeaders() map[string][]string {
	return r.resp.Header
}

func (r *Response) GetHeader(name string) []string {
	headers := r.GetHeaders()
	for k, v := range headers {
		if strings.EqualFold(name, k) {
			return v
		}
	}

	return nil
}

// GetHeaderLine 获取指定header的第一个值
func (r *Response) GetHeaderLine(name string) string {
	headers := r.GetHeader(name)
	if len(headers) > 0 {
		return headers[0]
	}
	return ""
}

func (r *Response) HasHeader(name string) bool {
	headers := r.GetHeaders()
	for k := range headers {
		if strings.EqualFold(name, k) {
			return true
		}
	}
	return false
}

func (r *Response) Err() error {
	return r.err
}

func (r *Response) Stream() chan []byte {
	return r.stream
}

func (r *Response) parseStream() {
	r.stream = make(chan []byte)
	decoder := eventsource.NewDecoder(r.resp.Body)

	go func() {
		defer r.resp.Body.Close()
		defer close(r.stream)

		// TODO 需要测试,不确定数据是否能全部接收并且响应
		for {
			event, err := decoder.Decode()
			if err != nil {
				if !errors.Is(err, io.EOF) {
					r.err = fmt.Errorf("read data failed: %v", err)
				}
				return
			}

			data := event.Data()
			// // TODO 此处不是很好，因为限制太强，只能用于特定的场景
			// if data == "[DONE]" {
			// 	// 消息读取完毕,成功退出
			// 	return
			// }

			r.stream <- []byte(data)
		}
	}()
}
