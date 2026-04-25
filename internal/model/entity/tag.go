package entity

// Tag 标签实体
type Tag struct {
	BaseEntity
	Name        string `gorm:"type:varchar(50);not null;index;comment:标签名称" json:"name"`
	Slug        string `gorm:"type:varchar(50);uniqueIndex;not null;comment:标签别名" json:"slug"`
	Description string `gorm:"type:varchar(200);comment:标签描述" json:"description"`
	PostCount   int    `gorm:"type:int unsigned;default:0;comment:文章数量" json:"post_count"`

	Articles []*Article `gorm:"many2many:article_tags" json:"articles,omitempty"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}
