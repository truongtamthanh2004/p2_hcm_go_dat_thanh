package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func TranslateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = w

		c.Next()

		if w.Status() >= 200 && w.Status() < 300 &&
			strings.Contains(w.Header().Get("Content-Type"), "application/json") {

			var resp map[string]interface{}
			if err := json.Unmarshal(w.body.Bytes(), &resp); err == nil {
				if msg, ok := resp["message"].(string); ok {
					if tFunc, exists := c.Get("T"); exists {
						if T, ok := tFunc.(func(string) string); ok {
							resp["message"] = T(msg)
						}
					}
				}

				orig := w.ResponseWriter

				for k, vv := range w.Header() {
					orig.Header()[k] = vv
				}

				orig.Header().Del("Content-Length")

				if ct := orig.Header().Get("Content-Type"); ct == "" {
					orig.Header().Set("Content-Type", "application/json; charset=utf-8")
				}

				status := w.Status()
				if status == 0 {
					status = http.StatusOK
				}
				if !w.Written() {
					orig.WriteHeader(status)
				}
				_ = json.NewEncoder(orig).Encode(resp)
				return
			}
		}

		orig := w.ResponseWriter
		for k, vv := range w.Header() {
			orig.Header()[k] = vv
		}
		orig.Header().Del("Content-Length")
		status := w.Status()
		if status == 0 {
			status = http.StatusOK
		}
		if !w.Written() {
			orig.WriteHeader(status)
		}
		_, _ = orig.Write(w.body.Bytes())
	}
}
