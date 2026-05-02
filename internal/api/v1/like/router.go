package like

import (
	"blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册点赞相关路由
func (ctrl *Controller) RegisterRoutes(r *gin.RouterGroup) {
	likes := r.Group("/likes")
	{
		// 获取文章点赞数（公开接口，无需登录）
		likes.GET("/article/:id/count", ctrl.GetArticleLikeCount)
		// 获取用户点赞状态（需要登录）
		likes.GET("/article/:id/status", middleware.Auth(), ctrl.GetUserLikeStatus)
		// 点赞文章（需要登录）
		likes.POST("/article/:id", middleware.Auth(), ctrl.LikeArticle)
		// 取消点赞（需要登录）
		likes.DELETE("/article/:id", middleware.Auth(), ctrl.UnlikeArticle)
		// 获取用户点赞的文章列表（需要登录）
		likes.GET("/user/articles", middleware.Auth(), ctrl.ListUserLikedArticles)
	}
}
