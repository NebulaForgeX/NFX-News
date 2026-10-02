package item

import "time"

type Item struct{ state State }
type State struct {
	ID, SourceID, OriginalID, Title, URL string
	MobileURL                            *string
	PubDate                              *time.Time
	Extra                                []byte
}

func NewFromState(st State) *Item { return &Item{state: st} }
func (i *Item) State() State      { return i.state }
