package middleware

import (
	"blog/internal/service/log_service"
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ResponseWriter 因为Writer没有read方法，继承Writer增加功能
type ResponseWriter struct {
	gin.ResponseWriter
	Body bytes.Buffer //临时保存返回体
	Head http.Header  //临时保存返回头
}

// Write 自己实现Write
func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.Body.Write(b) //保存返回体
	return w.ResponseWriter.Write(b)
}

// Header 重写header
func (w *ResponseWriter) Header() http.Header {
	return w.Head
}

// LogMiddleWare 日志绑定中间件
func LogMiddleWare(c *gin.Context) {
	//请求中间件
	log := log_service.NewActionLog(c)
	log.SetReqBody(c)
	c.Set("log", log)
	if c.Request.URL.Path == "/api/ai/article" {
		//因为流式输出会设置请求头
		c.Next()
		log.MiddlewareSave()
		return
	}
	//重写Writer方法,保存返回体
	res := &ResponseWriter{
		ResponseWriter: c.Writer,
		Head:           make(http.Header), //避免空指针
	}
	c.Writer = res
	c.Next()
	//响应中间件
	log.SetResBody(res.Body.Bytes())
	log.SetResHeader(res.Head)
	log.MiddlewareSave()
}
