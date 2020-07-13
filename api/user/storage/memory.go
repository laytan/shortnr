package storage

// Memory stores users in memory
type Memory struct {
	Users []User
}

// Set adds the user to the store
func (m *Memory) Set(user User) bool {
	m.Users = append(m.Users, user)
	return true
}

// Get returns the user with the given id and if it exists
func (m Memory) Get(id uint) (User, bool) {
	for _, user := range m.Users {
		if user.ID == id {
			return user, true
		}
	}
	return User{}, false
}

// GetByEmail returns the user with the given email and if it exists
func (m Memory) GetByEmail(email string) (User, bool) {
	for _, user := range m.Users {
		if user.Email == email {
			return user, true
		}
	}
	return User{}, false
}
