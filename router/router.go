package router

import (
	"auth/auth"
	"auth/service"
	"github.com/gin-gonic/gin"
	"log"
)

func Router() *gin.Engine {

	client, err := fbauth.InitAuth()
	if err != nil {
		log.Fatalln("failed to init firebase auth", err)
	}

	router := gin.Default()

	router.Use(fbauth.AuthJWT(client))

	router.GET("/user/:id", service.GetUser)
	router.GET("/user", service.GetAllUser)
	router.POST("/newuser", service.CreateUser)
	router.PUT("/user/:id", service.UpdateUser)
	router.DELETE("/deleteuser/:id", service.DeleteUser)

	return router
}
