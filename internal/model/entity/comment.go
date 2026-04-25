package entity

// Comment 评论实体
type Comment struct {
	BaseEntity
	ArticleID   uint   `gorm:"not null;index;comment:文章ID" json:"article_id"`
	UserID      uint   `gorm:"not null;index;comment:用户ID" json:"user_id"`
	ParentID    uint   `gorm:"default:0;index;comment:父评论ID" json:"parent_id"`
	ReplyToID   uint   `gorm:"default:0;comment:回复的目标评论ID" json:"reply_to_id"`
	ReplyToName string `gorm:"type:varchar(100);comment:回复的目标用户名" json:"reply_to_name"`
	Content     string `gorm:"type:varchar(1000);not null;comment:评论内容" json:"content"`
	Status      int    `gorm:"type:tinyint;default:1;comment:状态:1正常,2隐藏" json:"status"`

	Article *Article `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}
