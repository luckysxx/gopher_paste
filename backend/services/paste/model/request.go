package model

type CreatePasteRequest struct {
	Content  string `json:"content" binding:"required"`
	Language string `json:"language" binding:"required"`
}
