package auth

import (
	"github.com/gin-gonic/gin"
	"nacos/configuration"
	"nacos/model"
	"nacos/router"
	"nacos/token"
	"net/http"
)

func Auth(context *gin.Context) {
	if configuration.Configuration.Nacos.Auth.Enabled {
		accessToken := model.BindQuery(context, &model.AccessToken{})
		if accessToken.AccessToken == "" {
			accessToken = model.BindHeader(context, &model.AccessToken{})
		}
		claims, err := token.Manager.ParseToken(accessToken.AccessToken)
		if err != nil {
			router.AccessDenied.Msg("token expired!").Write(context, http.StatusForbidden)
			context.Abort()
			return
		}
		context.Set("username", claims.Subject)
	}
	context.Next()
}
