package request

type ArticleListRequest struct {
	PageRequest
	Keyword    string `form:"keyword"`
	CategoryID *uint  `form:"category_id"`
	TagID      *uint  `form:"tag_id"`
}

type AdminArticleListRequest struct {
	PageRequest
	Status  *int   `form:"status"`
	Keyword string `form:"keyword"`
}

type ArticleUpsertRequest struct {
	Title      string   `json:"title" binding:"required,min=1,max=255"`
	Slug       string   `json:"slug" binding:"omitempty,max=255"`
	Summary    string   `json:"summary" binding:"omitempty"`
	Content    string   `json:"content" binding:"required"`
	CategoryID *uint    `json:"category_id"`
	Tags       []string `json:"tags"`
	CoverImage string   `json:"cover_image" binding:"omitempty,max=500"`
	Status     int      `json:"status"`
}
