package log_service

import (
	"blog/internal/global"
	"blog/internal/model"
	"blog/internal/model/enum"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	e "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"reflect"
	"strings"
)

// ActionLog 操作日志
type ActionLog struct {
	c             *gin.Context
	level         enum.LogLevel
	log           *model.LogModel
	title         string
	reqBody       []byte
	resBody       []byte
	resHeader     http.Header
	showRes       bool     //是否显示响应体
	showReq       bool     //是否显示请求体
	showResHeader bool     //是否显示响应头
	showReqHeader bool     //是否显示请求头
	itemList      []string //响应和请求中间的显示内容
	isMiddleWare  bool     //是否是中间件调用
}

// SetLevel 设置level
func (actionLog *ActionLog) SetLevel(level enum.LogLevel) {
	actionLog.level = level
	return
}

// SetTitle 设置title
func (actionLog *ActionLog) SetTitle(title string) {
	actionLog.title = title
	return
}

// SetReqBody 保存请求体
func (actionLog *ActionLog) SetReqBody(c *gin.Context) {
	bytedata, _ := io.ReadAll(c.Request.Body)
	//重新创建请求体，重置文件偏移量
	c.Request.Body = io.NopCloser(bytes.NewReader(bytedata))
	actionLog.reqBody = bytedata
}

// SetResBody 保存返回体
func (actionLog *ActionLog) SetResBody(data []byte) {
	actionLog.resBody = data
}

func (actionLog *ActionLog) SetShowRes() {
	actionLog.showRes = true
}
func (actionLog *ActionLog) SetShowReq() {
	actionLog.showReq = true
}

// setItem 设置content
func (actionlog *ActionLog) setItem(label string, value any, logLevel enum.LogLevel) {
	var content string
	//断言
	switch reflect.TypeOf(value).Kind() {
	case reflect.Map, reflect.Struct, reflect.Slice:
		res, _ := json.Marshal(value)
		content = string(res)
	default:
		content = fmt.Sprintf("%v", value)
	}
	actionlog.itemList = append(actionlog.itemList, fmt.Sprintf(
		"<div class=\"log_item info\">%s<div class=\"log_item_label\">%s</div><div class=\"log_item_content\">%s</div></div>",
		logLevel.String(),
		label,
		content))
}
func (actionlog *ActionLog) SetItem(label string, value any) {
	actionlog.setItem(label, value, enum.LogLevelInfo)
}
func (actionlog *ActionLog) SetItemInfo(label string, value any) {
	actionlog.setItem(label, value, enum.LogLevelInfo)
}
func (actionlog *ActionLog) SetItemWarn(label string, value any) {
	actionlog.setItem(label, value, enum.LogLevelWarn)
}
func (actionlog *ActionLog) SetItemError(label string, value any) {
	actionlog.setItem(label, value, enum.LogLevelError)
}

// SetLink 设置超链接
func (actionlog *ActionLog) SetLink(label string, href string) {
	actionlog.itemList = append(actionlog.itemList,
		fmt.Sprintf("<div class=\"log_item link\"><div class=\"log_item_label\">%s</div><div class=\"log_item_content\"><a href=\"%s\" target=\"_blank\">%s</a></div></div> ",
			label,
			href, href))
}

// SetImage 设置图片
func (actionlog *ActionLog) SetImage(src string) {
	actionlog.itemList = append(actionlog.itemList, fmt.Sprintf("<div class=\"log_image\"><img src=\"%s\" alt=\"\"></div>", src))
}

// ShowReqHeader 展示请求头
func (actionlog *ActionLog) ShowReqHeader() {
	actionlog.showReqHeader = true
}

// ShowResHeader 展示响应头
func (actionlog *ActionLog) ShowResHeader() {
	actionlog.showResHeader = true
}
func (actionlog *ActionLog) SetResHeader(header http.Header) {
	actionlog.resHeader = header
}

// SetError 设置错误保存堆栈
func (actionlog *ActionLog) SetError(label string, err error) {
	msg := e.WithStack(err)
	logrus.Errorf("%s %s", label, err.Error())
	actionlog.itemList = append(actionlog.itemList, fmt.Sprintf("<div class=\"log_error\"><div class=\"line\"><div class=\"label\">%s</div><div class=\"value\">%s</div><div class=\"type\">%T</div></div><div class=\"stack\">%+v</div></div>",
		label, err, err, msg))
}

// MiddlewareSave 中间件调用逻辑
// 中间件在请求“生命周期”结束时自动记录日志，
// 而“其他服务”可能在业务逻辑中主动记录日志。
func (actionLog *ActionLog) MiddlewareSave() {
	_saveLog, _ := actionLog.c.Get("saveLog")
	saveLog, _ := _saveLog.(bool)
	if !saveLog {
		//说明不是主动调用的GetLog， 不需要保存
		return
	}
	if actionLog.log == nil {
		// 创建
		actionLog.isMiddleWare = true
		actionLog.Save()
		return
	}
	// 在视图里面save过，属于更新
	//响应头
	if actionLog.showResHeader {
		byteData, _ := json.Marshal(actionLog.resHeader)
		actionLog.itemList = append(actionLog.itemList, fmt.Sprintf("<div class=\"log_response_header\"><pre class=\"log_json_body\">%s</pre></div>", string(byteData)))
	}

	//设置响应
	if actionLog.showRes {
		actionLog.itemList = append(actionLog.itemList,
			fmt.Sprintf("<div class=\"log_response\"><pre class=\"log_json_body\">%s</pre></div>", string(actionLog.resBody)))
	}
	//收尾
	actionLog.Save()
}

// Save 保存操作到日志表 “业务代码”调用逻辑
// 要考虑到 业务代码主动调用和中间件调用的情况
func (actionLog *ActionLog) Save() (id uint) {
	if actionLog.log != nil {
		//说明之前已经保存过，不用创建，更新即可
		newContent := strings.Join(actionLog.itemList, "\n")
		global.DB.Model(actionLog.log).Updates(map[string]interface{}{
			"content": newContent,
		})
		actionLog.itemList = []string{}
		return actionLog.log.ID
	}
	//创建一个新的itemList，因为之前list里可能已经有内容
	var NewItemList []string
	// 请求头
	if actionLog.showReqHeader {
		byteData, _ := json.Marshal(actionLog.c.Request.Header)
		NewItemList = append(NewItemList, fmt.Sprintf("<div class=\"log_request_header\"><pre class=\"log_json_body\">%s</pre></div>", string(byteData)))
	}
	//设置请求
	if actionLog.showReq {
		//默认json
		NewItemList = append(NewItemList, fmt.Sprintf(
			"‹div class=\"log_request\"><div class=\"log_request_head\">«span class=\"log_request_method %s\">%s</span>‹span class=\"log_request_path\">%s</span></div>‹div class=\"log_request_body\">‹pre class=\"log_json_body\">%s</pre></div></div>", actionLog.c.Request.Method,
			strings.ToLower(actionLog.c.Request.Method),
			actionLog.c.Request.URL.String(),
			string(actionLog.reqBody)))
	}
	//设置中间的content
	NewItemList = append(NewItemList, actionLog.itemList...)

	//是中间件才拿响应
	if actionLog.isMiddleWare {
		//响应头
		if actionLog.showResHeader {
			byteData, _ := json.Marshal(actionLog.resHeader)
			NewItemList = append(NewItemList, fmt.Sprintf("<div class=\"log_response_header\"><pre class=\"log_json_body\">%s</pre></div>", string(byteData)))
		}

		//设置响应
		if actionLog.showRes {
			NewItemList = append(NewItemList,
				fmt.Sprintf("<div class=\"log_response\"><pre class=\"log_json_body\">%s</pre></div>", string(actionLog.resBody)))
		}
	}

	ip := actionLog.c.ClientIP()
	//TODO:通过jwt获取username
	//token := c.GetHeader("token")
	UserID := uint(1)
	log := model.LogModel{
		Type:    enum.LogActionTypes,
		Title:   actionLog.title,
		Content: strings.Join(NewItemList, "\n"),
		Level:   actionLog.level,
		UserID:  UserID,
		IP:      ip,
	}
	err := global.DB.Create(&log).Error
	if err != nil {
		logrus.Errorf("保存操作日志失败,%v", err)
		return
	}
	actionLog.log = &log
	return log.ID
}

// NewActionLog 创建新的Aclog对象
func NewActionLog(c *gin.Context) *ActionLog {
	return &ActionLog{
		c: c,
	}
}

// GetLog 获取创建的log对象
func GetLog(c *gin.Context) *ActionLog {
	_log, ok := c.Get("log")
	if !ok {
		return NewActionLog(c)
	}
	log, ok := _log.(*ActionLog)
	if !ok {
		return NewActionLog(c)
	}
	//如果主动调用GetLog
	c.Set("saveLog", true)
	return log
}
