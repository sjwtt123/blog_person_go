package entity

// Category 分类实体
type Category struct {
	BaseEntity
	Name        string `gorm:"type:varchar(100);not null;comment:分类名称" json:"name"`
	Slug        string `gorm:"type:varchar(100);uniqueIndex;not null;comment:分类别名" json:"slug"`
	Description string `gorm:"type:varchar(500);comment:分类描述" json:"description"`
	SortOrder   int    `gorm:"default:0;index;comment:排序" json:"sort_order"`
	ParentID    uint   `gorm:"default:0;index;comment:父分类ID" json:"parent_id"`
	PostCount   int    `gorm:"type:int unsigned;default:0;comment:文章数量" json:"post_count"`

	Articles []*Article `gorm:"foreignKey:CategoryID" json:"articles,omitempty"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}
