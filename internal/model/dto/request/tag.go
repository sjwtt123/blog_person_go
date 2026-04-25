package request

type TagUpsertRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Slug        string `json:"slug" binding:"omitempty,max=50"`
	Description string `json:"description" binding:"omitempty,max=200"`
}
