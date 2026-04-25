package request

type CategoryUpsertRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Slug        string `json:"slug" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
	SortOrder   int    `json:"sort_order"`
	ParentID    uint   `json:"parent_id"`
}
