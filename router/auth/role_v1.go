package auth

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nacos/consts"
	db "nacos/database"
	"nacos/model"
	"nacos/router"
	"nacos/util"
	"net/http"
)

func searchRole(context *gin.Context) {
	param := model.Bind(context, &model.SearchRole{})
	page := model.Bind(context, &model.Page{})
	var conditions []any
	blur := param.SearchType == "blur"
	if sql, arg := router.BlurQuery("username", param.Username, blur); arg != "" {
		conditions = append(conditions, []any{sql, arg})
	}
	if sql, arg := router.BlurQuery("role", param.Role, blur); arg != "" {
		conditions = append(conditions, []any{sql, arg})
	}
	context.JSON(http.StatusOK, model.PaginateResult[model.Role, model.RoleDetail](conditions, page))
}

func searchRoleName(context *gin.Context) {
	context.JSON(http.StatusOK, []string{})
}

func addRole(context *gin.Context) {
	param := model.Bind(context, &model.RoleInfo{})
	if param.Username == consts.DefaultUsername {
		context.String(http.StatusBadRequest, "cannot modify role of user %s", consts.DefaultUsername)
		return
	}
	role := &model.Role{}
	util.Copy(param, role)
	if err := db.GORM.Create(role).Error; errors.Is(err, gorm.ErrDuplicatedKey) {
		context.String(http.StatusBadRequest, "user %s already has role %s!", param.Username, param.Role)
	} else {
		router.OK.With("add role ok!").Ok(context)
	}
}

func deleteRole(context *gin.Context) {
	roleInfo := model.Bind(context, &model.RoleInfo{})
	username := roleInfo.Username
	if username == consts.DefaultUsername {
		context.String(http.StatusBadRequest, "cannot delete role of user %s", consts.DefaultUsername)
	} else {
		db.GORM.Delete(model.Role{}, roleInfo)
		router.OK.With(fmt.Sprintf("delete role of user %s ok!", username)).Ok(context)
	}
}
