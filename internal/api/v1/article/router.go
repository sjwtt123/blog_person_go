package article

import (
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/articles", ctrl.ListPublic)
	r.GET("/articles/timelines", ctrl.GetTimelineItems)
	r.GET("/articles/:id", ctrl.GetDetail)
	// 访问量统计
	r.PUT("/articles/:id/view", ctrl.IncrementViewCount)

	// 上传图片
	authGroup := r.Group("/user", middleware.Auth())
	{
		authGroup.POST("/upload/image", ctrl.UploadImage)
	}

	adminGroup := r.Group("/admin", middleware.Auth(), middleware.Admin(ctrl.userService))
	{
		adminGroup.POST("/articles", ctrl.Create)
		adminGroup.PUT("/articles/:id", ctrl.Update)
		adminGroup.DELETE("/articles/:id", ctrl.Delete)
		adminGroup.GET("/articles", ctrl.ListAdmin)
		adminGroup.POST("/clean/images", ctrl.CleanUnusedImages)
	}

}
