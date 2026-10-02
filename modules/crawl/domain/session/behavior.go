package session

func (s *Session) IsFinished() bool { return s.state.FinishedAt != nil }
func (s *Session) Status() string   { return s.state.Status }
