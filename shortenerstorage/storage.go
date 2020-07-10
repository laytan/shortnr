package shortenerstorage

// Storage defines the methods required to interact with a link storing method
type Storage interface {
	Get(id string) (string, bool)
	Set(id string, url string)
	Contains(url string) bool
}
