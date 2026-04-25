package middleware

import (
	"blog/internal/service"
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

func Admin(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		if userID == 0 {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response.BizError(c, err)
			c.Abort()
			return
		}

		if user.Role != "admin" {
			response.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}
