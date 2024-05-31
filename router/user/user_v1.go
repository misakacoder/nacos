package user

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nacos/database"
	"nacos/model"
	"nacos/router"
	"nacos/token"
	"nacos/util"
	"net/http"
	"time"
)

func RegisterV1(engine *gin.Engine) {
	auth := engine.Group(router.ApiV1 + "/auth")
	{
		auth.POST("/users/login", login)
	}
}

func login(context *gin.Context) {
	loginUser := model.Bind(context, &model.LoginUser{})
	username := loginUser.Username
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
	context.JSON(http.StatusOK, model.Token{AccessToken: tokenString, Ttl: claims.ExpiresAt - time.Now().Unix(), GlobalAdmin: true, Username: username})
}

func accessDenied(context *gin.Context, message string) {
	router.AccessDenied.Msg(message).Write(context, http.StatusForbidden)
}
