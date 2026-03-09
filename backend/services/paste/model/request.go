package model

type CreatePasteRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Language   string `json:"language" binding:"required"`
	Visibility string `json:"visibility"`
}

type UpdatePasteRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Language   string `json:"language" binding:"required"`
	Visibility string `json:"visibility"`
}
