package domain

type CategoryCreateRequest struct {
	Tag string `json:"tag"`
}

type CategoryUpdateRequest struct {
	ID  int64  `json:"id"`
	Tag string `json:"tag"`
}

type Category struct {
	ID    int64  `json:"id"`
	Tag   string `json:"tag"`
	Total int    `json:"total"`
}
