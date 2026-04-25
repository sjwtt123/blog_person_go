package request

// PageRequest 分页请求参数
type PageRequest struct {
	Page int `form:"page" binding:"required,min=1"`         // 页码，从1开始
	Size int `form:"size" binding:"required,min=1,max=100"` // 每页大小，最大100
}
