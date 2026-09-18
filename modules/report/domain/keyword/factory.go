package keyword

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

func New(accountID, profileID *string, group, word, kind string, limit int) (*Keyword, error) {
	word = strings.TrimSpace(word)
	kind = strings.ToLower(strings.TrimSpace(kind))
	if strings.HasPrefix(word, "!") {
		kind = "exclude"
		word = strings.TrimPrefix(word, "!")
	} else if strings.HasPrefix(word, "+") {
		kind = "required"
		word = strings.TrimPrefix(word, "+")
	}
	if kind == "" {
		kind = "include"
	}
	if group == "" {
		group = "default"
	}
	if err := validateWord(word, kind); err != nil {
		return nil, err
	}
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return NewFromState(State{
		ID: id, AccountID: accountID, ProfileID: profileID, GroupName: group, Word: word, Kind: kind, CountLimit: limit, CreatedAt: time.Now().UTC(),
	}), nil
}
