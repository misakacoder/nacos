package auth

import (
	"github.com/gin-gonic/gin"
	"nacos/router"
)

func RegisterV1(engine *gin.Engine) {
	auth := engine.Group(router.ApiV1 + "/auth")
	auth.POST("/users/login", login)

	auth.Use(Auth)
	{
		auth.GET("/users", searchUser)
		auth.GET("/users/search", searchUsername)
		auth.POST("/users", addUser)
		auth.PUT("/users", updateUser)
		auth.DELETE("/users", deleteUser)
	}
	{
		auth.GET("/roles", searchRole)
		auth.GET("/roles/search", searchRoleName)
		auth.POST("/roles", addRole)
		auth.DELETE("/roles", deleteRole)
	}
	{
		auth.GET("/permissions", searchPermission)
		auth.POST("/permissions", addPermission)
		auth.DELETE("/permissions", deletePermission)
	}
}
