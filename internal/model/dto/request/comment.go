package request

type CreateCommentRequest struct {
	Content     string `json:"content" binding:"required,min=1,max=1000"`
	ParentID    uint   `json:"parent_id"`
	ReplyToID   uint   `json:"reply_to_id"`
	ReplyToName string `json:"reply_to_name"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

type CommentListRequest struct {
	PageRequest
}

type AdminCommentListRequest struct {
	PageRequest
	ArticleID *uint  `form:"article_id"`
	Keyword   string `form:"keyword"`
}
