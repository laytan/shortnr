package shortenerstorage

type Storage interface {
	Get(id string) (string, bool)
	Set(id string, url string)
	Contains(url string) bool
}
