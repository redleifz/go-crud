package routes

import (
	"go-crud/controllers"
	"go-crud/middleware"

	"github.com/gin-gonic/gin"
)

// controllers.GetAllUsers
func UserRoute(apiGroup *gin.RouterGroup) {
	// Apply the middleware to the specific route
	apiGroup.GET("user/", middleware.VerifyTokenMiddleware(), middleware.CheckUserRole, controllers.GetAllUsers)

	apiGroup.GET("user/getexcel/", controllers.GetAllUsersExcelFile)

	apiGroup.POST("user/", controllers.CreateUser)
	apiGroup.POST("user/login/", controllers.UserLogin)

}
