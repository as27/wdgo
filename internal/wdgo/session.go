package wdgo

import (
	"fmt"
	"time"
)

// Session is used to track the time, when working on a card
type Session struct {
	EventSource
	Card  *Card
	Title string
	Start time.Time
	End   time.Time
}

func (s *Session) Event(cmd Cmd) error {
	switch cmd.Action {
	case "Start":
		s.Start = time.Now()
	case "End":
		s.End = time.Now()
	default:
		return fmt.Errorf("%w: %s", ErrNoSuchAction, cmd.Action)
	}
	return nil
}

func (s *Session) Find(id string) (Eventer, error) {
	if s.ID == id {
		return s, nil
	}
	return nil, ErrIDNotFound
}
