package comment

import (
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/articles/:id/comments", ctrl.List)

	authGroup := r.Group("/user", middleware.Auth())
	{
		authGroup.POST("/articles/:id/comments", ctrl.Create)
		authGroup.PUT("/comments/:id", ctrl.Update)
		authGroup.DELETE("/comments/:id", ctrl.Delete)
	}

	adminGroup := r.Group("/admin", middleware.Auth(), middleware.Admin(ctrl.userService))
	{
		adminGroup.GET("/comments", ctrl.ListAllAdmin)
		adminGroup.DELETE("/comments/:id", ctrl.AdminDelete)
	}
}
