package link

// Link is a link
type Link struct {
	ID        string `json:"id"`
	URL       string `json:"url" validate:"required,url,max=255"`
	UserID    uint   `json:"user_id" validate:"required"`
	CreatedAt string `json:"created_at"`
}
