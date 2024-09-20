package respx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yu1ec/go-pkg/errorx"
)

// ResponseWriter 是一个接口，可以同时被 gin.Context 和 http.ResponseWriter 实现
type ResponseWriter interface {
	Status(code int)
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

// GinResponseWriter 包装 gin.Context 以实现 ResponseWriter 接口
type GinResponseWriter struct {
	*gin.Context
}

func (g *GinResponseWriter) Status(code int) {
	g.Context.Status(code)
}

func (g *GinResponseWriter) Header() http.Header {
	return g.Context.Writer.Header()
}

func (g *GinResponseWriter) Write(data []byte) (int, error) {
	return g.Context.Writer.Write(data)
}

func (w *GinResponseWriter) WriteHeader(statusCode int) {
	w.Context.Writer.WriteHeader(statusCode)
}

type StandardResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (s *StandardResponseWriter) Status(code int) {
	s.statusCode = code
	s.ResponseWriter.WriteHeader(code)
}

// NewResponseWriter 创建一个新的 ResponseWriter
func NewResponseWriter(w any) ResponseWriter {
	switch v := w.(type) {
	case *gin.Context:
		return &GinResponseWriter{v}
	case http.ResponseWriter:
		return &StandardResponseWriter{v, 0}
	default:
		panic("unsupported writer type")
	}
}

// NoContent 是空响应
func NoContent(w ResponseWriter) {
	w.Status(http.StatusNoContent)
}

// PlainContent 是纯文本响应
func PlainContent(w ResponseWriter, data any) {
	w.Header().Set("Content-Type", "text/plain")
	w.Status(http.StatusOK)
	w.Write([]byte(data.(string)))
}

// JsonContent 是 JSON 响应
func JsonContent(w ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.Status(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// Pagination 是分页响应
func JsonPagination(w ResponseWriter, data any, total int64) {
	w.Header().Set("Content-Type", "application/json")
	w.Status(http.StatusOK)
	json.NewEncoder(w).Encode(gin.H{
		"data":  data,
		"total": total,
	})
}

// JsonResponseWithError 携带错误信息的Json响应
func JsonResponseWithError(w ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	if e, ok := err.(*errorx.Error); ok {
		w.Status(e.HttpStatusCode())
		json.NewEncoder(w).Encode(e.Data())
	} else {
		w.Status(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(gin.H{"error": err.Error()})
	}
}

type EventSourceMessage struct {
	Event string
	Data  string
	ID    string
}

type EventSourceHandler func(ctx context.Context, lastEventID string) (<-chan EventSourceMessage, error)

// EventSource Server-Sent Events 的实现
func EventSource(w ResponseWriter, r *http.Request, handler EventSourceHandler) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w.(http.ResponseWriter), "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	lastEventID := r.Header.Get("Last-Event-ID")
	events, err := handler(ctx, lastEventID)
	if err != nil {
		http.Error(w.(http.ResponseWriter), "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-events:
			if !ok {
				return
			}

			if msg.Event != "" {
				fmt.Fprintf(w, "event: %s\n", msg.Event)
			}

			if msg.ID != "" {
				fmt.Fprintf(w, "id: %s\n", msg.ID)
			}
			fmt.Fprintf(w, "data: %s\n\n", msg.Data)
			flusher.Flush()
		case <-time.After(15 * time.Second):
			// 发送一个心跳保持连接
			fmt.Fprintf(w, ": heartbeat\n\n")
			flusher.Flush()
		}
	}

}
