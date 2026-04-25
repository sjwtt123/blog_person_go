package response

type CategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	PostCount   int    `json:"post_count,omitempty"`
	SortOrder   int    `json:"sort_order,omitempty"`
	ParentID    uint   `json:"parent_id,omitempty"`
}
