package model

type PasteResponse struct {
	ShortLink string `json:"short_link"`
	Content   string `json:"content"`
	Language  string `json:"language"`
	CreatedAt string `json:"created_at"`
}
