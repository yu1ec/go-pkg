package respx_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

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
		respx.ErrorJsonResponse(w, err)
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
		respx.ErrorJsonResponse(w, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		var result map[string]string
		json.Unmarshal(rec.Body.Bytes(), &result)
		assert.Equal(t, "标准错误", result["error"])
	})
}
