package response

type LikeResponse struct {
	LikeCount int  `json:"like_count"`
	IsLiked   bool `json:"is_liked"`
}
