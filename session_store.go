package main

// SessionStore - the global session storage
type SessionStore map[string]*Session

// NewSessionStore - Initialize new session storage
func NewSessionStore() SessionStore {
	return SessionStore{}
}

// Acquire - get/create the session of the specified id
func (s SessionStore) Acquire(id string) *Session {
	if sess, exists := s[id]; exists {
		return sess
	}
	s[id] = NewSession()
	return s[id]
}

// Forget - forget the specified session
func (s SessionStore) Forget(id string) {
	delete(s, id)
}
