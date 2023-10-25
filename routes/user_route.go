package routes

import (
	"go-crud/controllers"

	"github.com/gin-gonic/gin"
)

// controllers.GetAllUsers
func UserRoute(apiGroup *gin.RouterGroup) {
	apiGroup.GET("user/", controllers.GetAllUsers)
	apiGroup.GET("user/getexcel/", controllers.GetAllUsersExcelFile)

	apiGroup.POST("user/", controllers.CreateUser)
	apiGroup.POST("user/login/", controllers.UserLogin)

	

}
