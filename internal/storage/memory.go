package storage

// MemoryStorage is a in memory link storage mainly used in testing
type MemoryStorage struct {
	InternalMap map[string]string
}

// Get returns the url associated with the given id
func (m MemoryStorage) Get(id string) (string, bool) {
	redirectURL, ok := m.InternalMap[id]
	return redirectURL, ok
}

// Set adds the given id to url to the map
func (m MemoryStorage) Set(id string, url string) {
	m.InternalMap[id] = url
}

// Contains checks if the map has the url in it
func (m MemoryStorage) Contains(url string) bool {
	contains := false
	for _, v := range m.InternalMap {
		if v == url {
			contains = true
			break
		}
	}
	return contains
}