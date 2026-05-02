package entity

// Like 点赞实体
type Like struct {
	BaseEntity
	UserID     uint   `gorm:"not null;index;comment:用户ID" json:"user_id"`
	TargetID   uint   `gorm:"not null;index;comment:目标ID" json:"target_id"`
	TargetType string `gorm:"type:varchar(20);not null;index;comment:目标类型:article,comment" json:"target_type"`
}

// TableName 指定表名
func (Like) TableName() string {
	return "likes"
}
