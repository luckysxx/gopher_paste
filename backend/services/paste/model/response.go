package model

type PasteResponse struct {
	ID         int64  `json:"id"`
	OwnerID    int64  `json:"owner_id"`
	Title      string `json:"title"`
	ShortLink  string `json:"short_link,omitempty"`
	Content    string `json:"content"`
	Language   string `json:"language"`
	Visibility string `json:"visibility"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
