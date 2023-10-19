package routes

import (
	"go-crud/controllers"

	"github.com/gin-gonic/gin"
)

// controllers.GetAllUsers
func UserRoute(apiGroup *gin.RouterGroup) {
	apiGroup.GET("user/", controllers.GetAllUsers)
	apiGroup.POST("user/", controllers.CreateUser)

}
