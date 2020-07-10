package link

import "net/url"

// IsLink checks if link is valid url
func IsLink(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
