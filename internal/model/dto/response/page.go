package response

import "math"

// PageResponse 分页响应结构
type PageResponse struct {
	List      interface{} `json:"list"`       // 数据列表
	Total     int64       `json:"total"`      // 总记录数
	Page      int         `json:"page"`       // 当前页码
	Size      int         `json:"size"`       // 每页大小
	TotalPage int         `json:"total_page"` // 总页数
}

// NewPageResponse 创建分页响应
func NewPageResponse(list interface{}, total int64, page, size int) *PageResponse {
	totalPage := int(math.Ceil(float64(total) / float64(size)))
	return &PageResponse{
		List:      list,
		Total:     total,
		Page:      page,
		Size:      size,
		TotalPage: totalPage,
	}
}
