package click

// Click is a click on a redirect/link
type Click struct {
	ID        uint   `json:"id"`
	LinkID    string `json:"link_id"`
	CreatedAt string `json:"created_at"`
}
