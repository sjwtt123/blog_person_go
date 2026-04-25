package response

import "time"

type AuthorResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar,omitempty"`
}

type ArticleListItemResponse struct {
	ID           uint              `json:"id"`
	Title        string            `json:"title"`
	Slug         string            `json:"slug,omitempty"`
	Summary      string            `json:"summary,omitempty"`
	CoverImage   string            `json:"cover_image,omitempty"`
	ViewCount    int               `json:"view_count"`
	CommentCount int               `json:"comment_count"`
	Author       *AuthorResponse   `json:"author,omitempty"`
	Category     *CategoryResponse `json:"category,omitempty"`
	Tags         []TagResponse     `json:"tags,omitempty"`
	PublishedAt  *time.Time        `json:"published_at,omitempty"`
	UpdatedAt    *time.Time        `json:"updated_at"`
	CreatedAt    *time.Time        `json:"created_at"`
	Status       int               `json:"status"`
}

type ArticleDetailResponse struct {
	ID           uint              `json:"id"`
	Title        string            `json:"title"`
	Slug         string            `json:"slug,omitempty"`
	Summary      string            `json:"summary,omitempty"`
	Content      string            `json:"content,omitempty"`
	CoverImage   string            `json:"cover_image,omitempty"`
	ViewCount    int               `json:"view_count"`
	CommentCount int               `json:"comment_count"`
	Author       *AuthorResponse   `json:"author,omitempty"`
	Category     *CategoryResponse `json:"category,omitempty"`
	Tags         []TagResponse     `json:"tags,omitempty"`
	PublishedAt  *time.Time        `json:"published_at,omitempty"`
	UpdatedAt    *time.Time        `json:"updated_at"`
	Status       int               `json:"status"`
}

type ArticleCreateResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug,omitempty"`
}

type TimelineItem struct {
	ID         uint              `json:"id"`
	Title      string            `json:"title"`
	Slug       string            `json:"slug"`
	CoverImage string            `json:"cover_image,omitempty"`
	ViewCount  int               `json:"view_count"`
	Category   *CategoryResponse `json:"category,omitempty"`
	CreatedAt  *time.Time        `json:"created_at"`
}
