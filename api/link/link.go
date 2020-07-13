package link

// Link is a link
type Link struct {
	URL string `validate:"required,url,max=255"`
}
