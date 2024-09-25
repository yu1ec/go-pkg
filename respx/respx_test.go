package respx_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yu1ec/go-pkg/errorx"
	"github.com/yu1ec/go-pkg/respx"
)

func TestNewResponseWriter(t *testing.T) {
	t.Run("Gin Context", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		w := respx.NewResponseWriter(c)
		assert.IsType(t, &respx.GinResponseWriter{}, w)
	})

	t.Run("HTTP ResponseWriter", func(t *testing.T) {
		rec := httptest.NewRecorder()
		w := respx.NewResponseWriter(rec)
		assert.IsType(t, &respx.StandardResponseWriter{}, w)
	})

	t.Run("Unsupported type", func(t *testing.T) {
		assert.Panics(t, func() {
			respx.NewResponseWriter("unsupported")
		})
	})
}

func TestNoContent(t *testing.T) {
	rec := httptest.NewRecorder()
	w := respx.NewResponseWriter(rec)
	respx.NoContent(w)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestPlainContent(t *testing.T) {
	rec := httptest.NewRecorder()
	w := respx.NewResponseWriter(rec)
	respx.PlainContent(w, "测试文本")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "text/plain", rec.Header().Get("Content-Type"))
	assert.Equal(t, "测试文本", rec.Body.String())
}

func TestJsonContent(t *testing.T) {
	rec := httptest.NewRecorder()
	w := respx.NewResponseWriter(rec)
	data := map[string]string{"key": "值"}
	respx.JsonContent(w, data)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	var result map[string]string
	json.Unmarshal(rec.Body.Bytes(), &result)
	assert.Equal(t, data, result)
}

func TestJsonPagination(t *testing.T) {
	rec := httptest.NewRecorder()
	w := respx.NewResponseWriter(rec)
	data := []any{"项目1", "项目2"}
	respx.JsonPagination(w, data, 10)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	var result map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &result)
	assert.Equal(t, data, result["data"])
	assert.Equal(t, float64(10), result["total"])
}

func TestErrorJsonResponse(t *testing.T) {
	t.Run("errorx.Error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		w := respx.NewResponseWriter(rec)
		err := errorx.NewError(http.StatusBadRequest, "000000", "无效的请求")
		respx.JsonResponseWithError(w, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		var result map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &result)
		assert.Equal(t, "无效的请求", result["reason"])
	})

	t.Run("standard error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		w := respx.NewResponseWriter(rec)
		err := errors.New("标准错误")
		respx.JsonResponseWithError(w, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		var result map[string]string
		json.Unmarshal(rec.Body.Bytes(), &result)
		assert.Equal(t, "标准错误", result["error"])
	})
}

// 模拟 ResponseWriter 和 Flusher
type mockResponseWriter struct {
	headers    http.Header
	body       strings.Builder
	statusCode int
	flushed    bool
}

func newMockResponseWriter() *mockResponseWriter {
	return &mockResponseWriter{
		headers: make(http.Header),
	}
}

func (m *mockResponseWriter) Header() http.Header {
	return m.headers
}

func (m *mockResponseWriter) Write(b []byte) (int, error) {
	return m.body.Write(b)
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}

func (m *mockResponseWriter) Status(code int) {
	m.statusCode = code
}

func (m *mockResponseWriter) Flush() {
	m.flushed = true
}

func TestEventSource(t *testing.T) {
	t.Run("正常情况", func(t *testing.T) {
		w := newMockResponseWriter()
		r := httptest.NewRequest("GET", "/events", nil)

		handler := func(ctx context.Context, lastEventID string) (<-chan respx.EventSourceMessage, error) {
			ch := make(chan respx.EventSourceMessage, 1)
			ch <- respx.EventSourceMessage{Event: "test", Data: "hello", ID: "1"}
			close(ch)
			return ch, nil
		}

		done := make(chan bool)
		go func() {
			respx.EventSource(w, r, handler)
			done <- true
		}()

		select {
		case <-done:
		case <-time.After(1 * time.Second):
			t.Fatal("EventSource 没有在预期时间内完成")
		}

		expectedOutput := "event: test\nid: 1\ndata: hello\n\n"
		if w.body.String() != expectedOutput {
			t.Errorf("输出不匹配。期望：%q，实际：%q", expectedOutput, w.body.String())
		}

		if !w.flushed {
			t.Error("Flush 未被调用")
		}
	})

	t.Run("不支持 Flusher", func(t *testing.T) {
		rec := httptest.NewRecorder()
		w := &respx.StandardResponseWriter{ResponseWriter: rec}
		r := httptest.NewRequest("GET", "/events", nil)

		handler := func(ctx context.Context, lastEventID string) (<-chan respx.EventSourceMessage, error) {
			return make(chan respx.EventSourceMessage), nil
		}

		respx.EventSource(w, r, handler)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("期望状态码 %d，实际：%d", http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestEvents(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("BeforeResponseEvent", func(t *testing.T) {
		respx.ClearEventHandlers()
		defer respx.ClearEventHandlers()

		eventFired := false
		respx.RegisterEvent(respx.BeforeResponse, func(w respx.ResponseWriter, data any) {
			eventFired = true
			assert.Equal(t, "Hello, Event!", data)
		})

		w := httptest.NewRecorder()
		respx.PlainContent(respx.NewResponseWriter(w), "Hello, Event!")
		assert.True(t, eventFired)
	})

	t.Run("AfterResponseEvent", func(t *testing.T) {
		respx.ClearEventHandlers()
		defer respx.ClearEventHandlers()

		eventFired := false
		respx.RegisterEvent(respx.AfterResponse, func(w respx.ResponseWriter, data any) {
			eventFired = true
			assert.Equal(t, "Hello, Event!", data)
		})

		w := httptest.NewRecorder()
		respx.PlainContent(respx.NewResponseWriter(w), "Hello, Event!")
		assert.True(t, eventFired)
	})

	t.Run("MultipleEvents", func(t *testing.T) {
		respx.ClearEventHandlers()
		defer respx.ClearEventHandlers()

		beforeFired := false
		afterFired := false

		respx.RegisterEvent(respx.BeforeResponse, func(w respx.ResponseWriter, data any) {
			beforeFired = true
			assert.Equal(t, map[string]string{"message": "Hello, Events!"}, data)
		})
		respx.RegisterEvent(respx.AfterResponse, func(w respx.ResponseWriter, data any) {
			afterFired = true
			assert.Equal(t, map[string]string{"message": "Hello, Events!"}, data)
		})

		w := httptest.NewRecorder()
		respx.JsonContent(respx.NewResponseWriter(w), map[string]string{"message": "Hello, Events!"})

		assert.True(t, beforeFired)
		assert.True(t, afterFired)
	})
}
