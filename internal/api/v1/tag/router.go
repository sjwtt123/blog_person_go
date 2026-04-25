package tag

import (
	"blog/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/tags", ctrl.ListAll)

	admin := r.Group("/admin", middleware.Auth(), middleware.Admin(ctrl.userService))
	{
		admin.POST("/tags", ctrl.Create)
		admin.PUT("/tags/:id", ctrl.Update)
		admin.DELETE("/tags/:id", ctrl.Delete)
	}
}
