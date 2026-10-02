package snapshot

func (s *Snapshot) Mode() string   { return s.state.Mode }
func (s *Snapshot) Title() string  { return s.state.Title }
func (s *Snapshot) ItemCount() int { return s.state.ItemCount }
