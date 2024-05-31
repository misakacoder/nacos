package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ApiV1   = "/nacos/v1"
	ApiV2   = "/nacos/v2"
	Success = Result{Code: 0, Message: "success", Meaning: "成功"}
	OK      = Result{Code: 200, Message: "success", Meaning: "成功"}
	Error   = Result{Code: 500, Message: "error", Meaning: "失败"}

	ParameterMissing = Result{Code: 10000, Message: "parameter missing", Meaning: "参数缺失"}
	AccessDenied     = Result{Code: 10001, Message: "access denied", Meaning: "访问拒绝"}

	ResourceNotFound    = Result{Code: 20004, Message: "resource not found", Meaning: "资源未找到"}
	ConfigListenerError = Result{Code: 20007, Message: "config listener error", Meaning: "监听配置错误"}

	ServerError = Result{Code: 30000, Message: "server error", Meaning: "其他内部错误"}
)

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Meaning string `json:"-"`
	Data    any    `json:"data,omitempty"`
}

func (result *Result) Msg(message string) *Result {
	return &Result{Code: result.Code, Message: message}
}

func (result *Result) With(v any) *Result {
	return &Result{Code: result.Code, Message: result.Message, Data: v}
}

func (result *Result) Ok(context *gin.Context) {
	context.JSON(http.StatusOK, result)
}

func (result *Result) Error(context *gin.Context) {
	context.JSON(http.StatusInternalServerError, result)
}

func (result *Result) Write(context *gin.Context, status int) {
	context.JSON(status, result)
}

func NotFound(context *gin.Context) {
	context.JSON(http.StatusNotFound, ResourceNotFound)
}
