package auth

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nacos/consts"
	"nacos/database"
	"nacos/model"
	"nacos/router"
	"nacos/token"
	"nacos/util"
	"net/http"
	"time"
)

func login(context *gin.Context) {
	loginUser := model.Bind(context, &model.UserInfo{})
	username := loginUser.Username.Username
	password := util.MD5(loginUser.Password)
	user := &model.User{Username: username}
	var message string
	if err := db.GORM.Where(user).First(user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		message = fmt.Sprintf("user %s not found", username)
	} else if password != user.Password {
		message = "password error"
	}
	if message != "" {
		accessDenied(context, message)
		return
	}
	tokenString, claims := token.Manager.CreateToken(username)
	var count int64
	db.GORM.Where(&model.Role{Username: username, Role: consts.DefaultRole}).Count(&count)
	isGlobalAdmin := util.ConditionalExpression(count > 0, true, false)
	context.JSON(http.StatusOK, model.Token{AccessToken: tokenString, Ttl: claims.ExpiresAt - time.Now().Unix(), GlobalAdmin: isGlobalAdmin, Username: username})
}

func searchUser(context *gin.Context) {
	param := model.Bind(context, &model.SearchUser{})
	page := model.Bind(context, &model.Page{})
	var conditions []any
	if sql, arg := router.BlurQuery("username", param.Username, param.SearchType == "blur"); arg != "" {
		conditions = append(conditions, []any{sql, arg})
	}
	context.JSON(http.StatusOK, model.PaginateResult[model.User, model.UserDetail](conditions, page))
}

func searchUsername(context *gin.Context) {
	username := model.Bind(context, &model.Username{})
	field := "username"
	sql, arg := router.BlurQuery(field, username.Username, true)
	var usernames []string
	db.GORM.Model(model.User{}).Select(field).Where(sql, arg).Find(&usernames)
	context.JSON(http.StatusOK, usernames)
}

func addUser(context *gin.Context) {
	userInfo := model.Bind(context, &model.UserInfo{})
	user := &model.User{Username: userInfo.Username.Username, Password: util.MD5(userInfo.Password)}
	if err := db.GORM.Create(user).Error; errors.Is(err, gorm.ErrDuplicatedKey) {
		context.String(http.StatusBadRequest, "user %s already exist!", userInfo.Username)
	} else {
		router.OK.With("create user ok!").Ok(context)
	}
}

func updateUser(context *gin.Context) {
	param := model.Bind(context, &model.UpdateUser{})
	db.GORM.Where(&model.User{Username: param.Username.Username}).Updates(&model.User{Password: util.MD5(param.NewPassword)})
	router.OK.With("update user ok!").Ok(context)
}

func deleteUser(context *gin.Context) {
	username := model.Bind(context, &model.Username{})
	if username.Username == consts.DefaultUsername {
		context.String(http.StatusBadRequest, "cannot delete admin: nacos")
	} else {
		db.GORM.Delete(model.User{}, username)
		router.OK.With("delete user ok!").Ok(context)
	}
}

func accessDenied(context *gin.Context, message string) {
	router.AccessDenied.Msg(message).Write(context, http.StatusForbidden)
}
