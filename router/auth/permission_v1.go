package auth

import (
	"github.com/gin-gonic/gin"
	"nacos/model"
	"nacos/router"
	"net/http"
)

func searchPermission(context *gin.Context) {
	context.JSON(http.StatusOK, model.PageResult{PageNum: 1, List: []struct{}{}})
}

func addPermission(context *gin.Context) {
	router.OK.Msg("add permission ok!").Ok(context)
}

func deletePermission(context *gin.Context) {
	router.OK.Msg("delete permission ok!").Ok(context)
}
