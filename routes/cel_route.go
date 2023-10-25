package routes

import (
	"go-crud/controllers"

	"github.com/gin-gonic/gin"
)

func CelRoutes(apiGroup *gin.RouterGroup) {
	apiGroup.GET("cel/", controllers.GetCel)
}
