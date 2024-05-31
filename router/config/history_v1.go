package config

import (
	"github.com/gin-gonic/gin"
	"nacos/database"
	"nacos/model"
	"net/http"
)

func listHistoryConfig(context *gin.Context) {
	namespaceID := context.DefaultQuery("tenant", "")
	configKey := model.ConfigKey{NamespaceID: &namespaceID}
	var configKeys []model.ConfigKey
	db.GORM.Model(model.HistoryConfigInfo{}).Where(&configKey).Find(&configKeys)
	context.JSON(http.StatusOK, configKeys)
}

func queryHistoryConfig(context *gin.Context) {
	if searchType := context.Query("search"); searchType != "" {
		searchHistoryConfig(context, searchType)
	} else {
		getHistoryConfig(context)
	}
}

func searchHistoryConfig(context *gin.Context, searchType string) {
	configKey := model.Bind(context, &model.ConfigKey{})
	configKey.SetNamespaceID()
	page := model.Bind(context, &model.Page{})
	page.OrderBy = "id desc"
	context.JSON(http.StatusOK, model.PaginateResult[model.HistoryConfigInfo, model.HistoryConfigDetail](configKey, page))
}

func getHistoryConfig(context *gin.Context) {
	param := model.Bind(context, &model.GetHistoryConfig{})
	param.SetNamespaceID()
	result := model.HistoryConfigDetail{}
	db.GORM.Model(model.HistoryConfigInfo{}).Where(param).Find(&result)
	context.JSON(http.StatusOK, result)
}
