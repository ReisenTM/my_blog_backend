package resp

import (
	"blog/internal/utils/validate"
	"github.com/gin-gonic/gin"
	"net/http"
)

var empty = map[string]any{}

type Code int
type Response struct {
	Code    Code   `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func (c Code) String() string {
	switch c {
	case SuccessCode:
		return "成功"
	case FailValidCode:
		return "失败"
	}
	return ""
}

const (
	SuccessCode   Code = 0
	FailValidCode Code = -1
)

func (r *Response) Json(c *gin.Context) {
	c.JSON(http.StatusOK, r)
}

func OK(data any, msg string, c *gin.Context) {
	response := Response{Code: SuccessCode, Message: msg, Data: data}
	response.Json(c)
}
func OKWithMsg(msg string, c *gin.Context) {
	response := Response{Code: SuccessCode, Message: msg, Data: empty}
	response.Json(c)
}
func OkWithData(data any, c *gin.Context) {
	resp := Response{SuccessCode, data, "成功"}
	resp.Json(c)
}
func OkWithList(list any, count int, c *gin.Context) {
	resp := Response{SuccessCode, map[string]any{
		"list":  list,
		"count": count,
	}, "成功"}
	resp.Json(c)
}
func FailWithCode(code Code, c *gin.Context) {
	response := Response{Code: code, Message: "失败", Data: empty}
	response.Json(c)
}
func FailWithMsg(msg string, c *gin.Context) {
	response := Response{Code: FailValidCode, Message: msg, Data: empty}
	response.Json(c)
}
func FailWithData(data any, c *gin.Context) {
	response := Response{
		Code:    SuccessCode,
		Message: "成功",
		Data:    data,
	}
	response.Json(c)
}
func FailWithError(err error, c *gin.Context) {
	data, msg := validate.ValidateError(err)
	response := Response{Code: FailValidCode, Message: msg, Data: data}
	response.Json(c)
}
