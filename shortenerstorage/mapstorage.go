package shortenerstorage

type MapStorage struct {
	InternalMap map[string]string
}

func (m MapStorage) Get(id string) (string, bool) {
	redirectURL, ok := m.InternalMap[id]
	return redirectURL, ok
}

func (m MapStorage) Set(id string, url string) {
	m.InternalMap[id] = url
}

func (m MapStorage) Contains(url string) bool {
	contains := false
	for _, v := range m.InternalMap {
		if v == url {
			contains = true
			break
		}
	}
	return contains
}
