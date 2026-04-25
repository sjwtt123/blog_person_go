package entity

import "time"

// Article 文章实体
type Article struct {
	BaseEntity
	Title        string     `gorm:"type:varchar(255);not null;comment:文章标题" json:"title"`
	Slug         string     `gorm:"type:varchar(255);index;comment:文章别名" json:"slug"`
	Summary      string     `gorm:"type:text;comment:文章摘要" json:"summary"`
	Content      string     `gorm:"type:longtext;not null;comment:文章内容" json:"content"`
	CoverImage   string     `gorm:"type:varchar(500);comment:封面图片" json:"cover_image"`
	ViewCount    int        `gorm:"type:int unsigned;default:0;comment:浏览次数" json:"view_count"`
	CommentCount int        `gorm:"type:int unsigned;default:0;comment:评论数" json:"comment_count"`
	Status       int        `gorm:"type:tinyint;default:1;index;comment:状态:1发布,2草稿,3隐藏" json:"status"`
	PublishedAt  *time.Time `gorm:"index;comment:发布时间" json:"published_at"`

	AuthorID uint  `gorm:"not null;index;comment:作者ID" json:"author_id"`
	Author   *User `gorm:"foreignKey:AuthorID" json:"author,omitempty"`

	Category   *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	CategoryID *uint     `gorm:"index;comment:分类ID" json:"category_id"`

	Tags []*Tag `gorm:"many2many:article_tags" json:"tags,omitempty"`
}

// TableName 指定表名
func (Article) TableName() string {
	return "articles"
}
