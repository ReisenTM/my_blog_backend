package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

type RequestContent struct {
	URL     *url.URL
	Headers http.Header
	Body    bytes.Buffer
}
type MyHeaders struct {
	http.Header
}

func (h *MyHeaders) String() string {
	b, _ := json.Marshal(h.Header)
	return string(b)
}

func (r *RequestContent) String() string {
	return fmt.Sprintf(
		"URL: %s\nHeaders: %v\nBody: %s\n",
		r.URL.String(),
		r.Headers,
		r.Body.String(),
	)
}

func LoggerMiddleware(c *gin.Context) {
	reqSaver := RequestContent{}
	reqSaver.URL = c.Request.URL
	reqSaver.Headers = c.Request.Header
	bodybytes, _ := io.ReadAll(c.Request.Body)
	reqSaver.Body.Write(bodybytes)
	logrus.Infof("reqSaver body: %s", reqSaver.String())
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodybytes)) //重置位置
	c.Next()
}
