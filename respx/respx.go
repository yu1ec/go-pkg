package respx

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yu1ec/go-pkg/errorx"
)

// ResponseWriter 是一个接口，可以同时被 gin.Context 和 http.ResponseWriter 实现
type ResponseWriter interface {
	Status(code int)
	Header() http.Header
	Write([]byte) (int, error)
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
