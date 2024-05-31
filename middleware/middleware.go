package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/misakacoder/logger"
	"nacos/router"
	"nacos/util"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func NetWork(context *gin.Context) {
	start := time.Now()
	context.Next()
	duration := time.Since(start)
	request := context.Request
	logger.Info("%s %s %s %d %dms", util.GetClientIP(context), request.Method, request.URL.Path, context.Writer.Status(), duration.Milliseconds())
}

func CSRF(context *gin.Context) {
	method := context.Request.Method
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, HEAD")
	context.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token")
	context.Header("Access-Control-Expose-Headers", "Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, Content-Length")
	context.Header("Access-Control-Allow-Credentials", "true")
	if method == "OPTIONS" {
		context.AbortWithStatus(http.StatusOK)
	}
	context.Next()
}

func Recovery(context *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			serverError := router.ServerError.With(nil)
			switch tp := err.(type) {
			case error:
				serverError.Data = tp.Error()
			case router.Result:
				serverError = &tp
			case *router.Result:
				serverError = tp
			case string:
				serverError.Data = tp
			default:
				serverError.Data = fmt.Sprintf("%v", tp)
			}
			logger.Error("%v", getStackTrace(err))
			serverError.Error(context)
			context.Abort()
		}
	}()
	context.Next()
}

func getStackTrace(err any) string {
	stackTrace := strings.Builder{}
	stackTrace.WriteString(fmt.Sprintf("%v", err))
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		stackTrace.WriteString(fmt.Sprintf("\n - %s:%d (0x%x)", file, line, pc))
	}
	return stackTrace.String()
}
