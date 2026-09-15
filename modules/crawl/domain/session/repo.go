package session

import "context"

type Repo struct {
	Create Create
	Update Update
}
type Create interface {
	New(ctx context.Context, s *Session) error
}
type Update interface {
	Generic(ctx context.Context, s *Session) error
}
