package namespace

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	db "nacos/database"
	"nacos/database/dbutil"
	"nacos/model"
	"nacos/router"
	"nacos/router/user"
	"nacos/util"
	"net/http"
)

func RegisterV1(engine *gin.Engine) {
	namespace := engine.Group(router.ApiV1+"/console/namespaces", user.Auth)
	{
		namespace.GET("", queryNamespace)
		namespace.POST("", saveNamespace)
		namespace.PUT("", updateNamespace)
		namespace.DELETE("", deleteNamespace)
	}
}

func queryNamespace(context *gin.Context) {
	if context.Query("show") == "all" {
		getNamespace(context)
	} else if context.Query("checkNamespaceIdExist") == "true" {
		checkNamespace(context)
	} else {
		listNamespace(context)
	}
}

func getNamespace(context *gin.Context) {
	searchNamespace := model.Bind(context, &model.SearchNamespace{})
	searchNamespace.SetNamespaceID()
	namespaceInfo := model.NamespaceInfo{}
	if err := db.GORM.Where(searchNamespace).First(&namespaceInfo).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		context.String(http.StatusInternalServerError, "namespaceId %s not exist", searchNamespace.NamespaceID)
	} else {
		context.JSON(http.StatusOK, getNamespaceInfoVO(&namespaceInfo))
	}
}

func checkNamespace(context *gin.Context) {
	customNamespaceId := context.Query("customNamespaceId")
	var count int64
	if customNamespaceId != "" {
		db.GORM.Where(model.NamespaceInfo{NamespaceID: customNamespaceId}).Count(&count)
	}
	context.String(http.StatusOK, "%v", count > 0)
}

func listNamespace(context *gin.Context) {
	var namespaceInfos []model.NamespaceInfo
	var namespaceInfoVOs []model.NamespaceDetail
	db.GORM.Find(&namespaceInfos)
	for _, namespaceInfo := range namespaceInfos {
		namespaceInfoVOs = append(namespaceInfoVOs, getNamespaceInfoVO(&namespaceInfo))
	}
	context.JSON(http.StatusOK, router.OK.With(namespaceInfoVOs))
}

func saveNamespace(context *gin.Context) {
	addNamespace := model.Bind(context, &model.AddNamespace{})
	data := &model.NamespaceInfo{}
	util.Copy(addNamespace, &data)
	if data.NamespaceID == "" {
		namespaceID, _ := uuid.NewRandom()
		data.NamespaceID = namespaceID.String()
	}
	err := db.GORM.Create(data).Error
	context.String(http.StatusOK, "%v", err == nil)
}

func updateNamespace(context *gin.Context) {
	param := model.Bind(context, &model.UpdateNamespace{})
	data := model.NamespaceInfo{}
	util.Copy(&param, &data)
	db.GORM.Model(model.NamespaceInfo{}).Where(&model.UpdateNamespace{NamespaceID: param.NamespaceID}).Updates(data)
	context.String(http.StatusOK, "%v", true)
}

func deleteNamespace(context *gin.Context) {
	searchNamespace := model.Bind(context, &model.SearchNamespace{})
	namespaceID := searchNamespace.NamespaceID
	if *namespaceID == "" {
		router.ServerError.With("deletePublic").Ok(context)
		return
	}
	configCount := countConfig(*namespaceID)
	if configCount == 0 {
		db.Transaction(func(tx *gorm.DB) {
			dbutil.PanicError(tx.Delete(model.HistoryConfigInfo{}, searchNamespace))
			affected := dbutil.PanicError(tx.Delete(model.NamespaceInfo{}, searchNamespace)).RowsAffected
			context.String(http.StatusOK, "%v", util.ConditionalExpression(affected > 0, true, false))
		})
	} else {
		router.ServerError.With("existConfiguration").Ok(context)
	}
}

func getNamespaceInfoVO(namespaceInfo *model.NamespaceInfo) model.NamespaceDetail {
	namespaceID := namespaceInfo.NamespaceID
	vo := &model.NamespaceDetail{}
	util.Copy(namespaceInfo, vo)
	vo.Quota = 200
	vo.Type = util.ConditionalExpression(namespaceID == "", 0, 2)
	vo.ConfigCount = countConfig(namespaceID)
	return *vo
}

func countConfig(namespaceID string) int {
	var count int64
	db.GORM.Model(model.ConfigInfo{}).Where(model.ConfigKey{NamespaceID: &namespaceID}).Count(&count)
	return int(count)
}
