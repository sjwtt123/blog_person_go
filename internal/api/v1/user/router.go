package user

import (
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册用户路由
func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/about", ctrl.GetAbout)

	userGroup := r.Group("/user")
	userGroup.Use(middleware.Auth())
	{
		userGroup.GET("/profile", ctrl.GetProfile)
		userGroup.PUT("/profile", ctrl.UpdateProfile)
		userGroup.POST("/password", ctrl.ChangePassword)
	}
	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.Auth(), middleware.Admin(ctrl.userService))
	{
		adminGroup.GET("/users", ctrl.ListUsers)
		adminGroup.POST("/users", ctrl.CreateUser)
		adminGroup.PUT("/users/:id", ctrl.UpdateUser)
		adminGroup.DELETE("/users/:id", ctrl.DeleteUser)
		adminGroup.PUT("/about", ctrl.UpdateAbout)
	}
}
