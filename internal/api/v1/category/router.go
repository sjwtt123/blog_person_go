package category

import (
	"blog/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/categories", ctrl.ListAll)

	admin := r.Group("/admin", middleware.Auth(), middleware.Admin(ctrl.userService))
	{
		admin.POST("/categories", ctrl.Create)
		admin.PUT("/categories/:id", ctrl.Update)
		admin.DELETE("/categories/:id", ctrl.Delete)
	}
}
