package api

import (
	"blog/internal/api/v1/article"
	"blog/internal/api/v1/auth"
	"blog/internal/api/v1/category"
	"blog/internal/api/v1/comment"
	"blog/internal/api/v1/like"
	"blog/internal/api/v1/tag"
	"blog/internal/api/v1/user"
	"blog/internal/middleware"
	"blog/internal/service"

	"github.com/gin-gonic/gin"
)

// Router 路由
type Router struct {
	userCtrl     *user.Controller
	authCtrl     *auth.Controller
	articleCtrl  *article.Controller
	categoryCtrl *category.Controller
	tagCtrl      *tag.Controller
	commentCtrl  *comment.Controller
	likeCtrl     *like.Controller
}

// NewRouter 创建路由
func NewRouter(
	userService service.UserService,
	authService service.AuthService,
	articleService service.ArticleService,
	categoryService service.CategoryService,
	tagService service.TagService,
	viewCountService service.ViewCountService,
	commentService service.CommentService,
	likeService service.LikeService,
	uploadCleanService service.UploadCleanService,
	uploadPath string,
) *Router {
	return &Router{
		userCtrl:     user.NewController(userService),
		authCtrl:     auth.NewController(authService, userService),
		articleCtrl:  article.NewController(articleService, userService, viewCountService, uploadCleanService, uploadPath),
		categoryCtrl: category.NewController(categoryService, userService),
		tagCtrl:      tag.NewController(tagService, userService),
		commentCtrl:  comment.NewController(commentService, userService),
		likeCtrl:     like.NewController(likeService),
	}
}

// Setup 设置路由
func (r *Router) Setup(engine *gin.Engine) {
	// 全局中间件
	engine.Use(middleware.Recovery())
	engine.Use(middleware.Logger())
	engine.Use(middleware.CORS())

	// 静态资源服务 - 提供上传文件访问
	engine.Static("/uploads", "./uploads")

	// 健康检查
	engine.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "blog API is running",
		})
	})

	// API v1 路由组
	v1 := engine.Group("/api/v1")
	{
		// 认证路由
		r.authCtrl.RegisterRoutes(v1)

		// 用户路由
		r.userCtrl.RegisterRoutes(v1)

		// 文章路由
		r.articleCtrl.RegisterRoutes(v1)

		//分类路由
		r.categoryCtrl.RegisterRoutes(v1)

		//标签路由
		r.tagCtrl.RegisterRoutes(v1)

		//评论路由
		r.commentCtrl.RegisterRoutes(v1)

		//点赞路由
		r.likeCtrl.RegisterRoutes(v1)
	}
}
