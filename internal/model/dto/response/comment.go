package response

import "time"

type CommentResponse struct {
	ID          uint               `json:"id"`
	ArticleID   uint               `json:"article_id"`
	UserID      uint               `json:"user_id"`
	ParentID    uint               `json:"parent_id"`
	ReplyToID   uint               `json:"reply_to_id"`
	ReplyToName string             `json:"reply_to_name"`
	Content     string             `json:"content"`
	Status      int                `json:"status"`
	CreatedAt   *time.Time         `json:"created_at"`
	UpdatedAt   *time.Time         `json:"updated_at"`
	User        *CommentUser       `json:"user,omitempty"`
	Article     *CommentArticle    `json:"article,omitempty"`
	Replies     []*CommentResponse `json:"replies,omitempty"`
}

type CommentUser struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar,omitempty"`
	Role     string `json:"role"`
}

type CommentArticle struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}
