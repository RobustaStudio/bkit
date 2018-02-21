package main

// Session - contains the session related data
type Session struct {
	CurrentContainer   string
	CurrentElement     string
	ElementsHistory    string
	ExpectingUserInput bool
	EOF                bool
	Data               map[string]interface{}
}

// NewSession - Initialize a new session
func NewSession() *Session {
	return &Session{
		Data: map[string]interface{}{},
	}
}

// Clear - Clear the current initialized session
func (s *Session) Clear() {
	s.CurrentElement = ""
	s.CurrentContainer = ""
	s.ElementsHistory = ""
	s.ExpectingUserInput = false
	s.EOF = false
	s.Data = map[string]interface{}{}
}
