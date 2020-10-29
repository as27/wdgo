package wdgo

import (
	"fmt"
	"time"
)

// Session is used to track the time, when working on a card
type Session struct {
	EventSource
	Card  *Card
	Start time.Time
	End   time.Time
}

func (s *Session) Event(cmd Cmd) error {
	var err error
	switch cmd.Action {
	case "Name":
		s.Name = cmd.Value
	case "Start":
		s.Start, err = time.Parse(TimeFormat, cmd.Value)
		if err != nil {
			return fmt.Errorf("Session.Start: %w", err)
		}
	case "End":
		s.End, err = time.Parse(TimeFormat, cmd.Value)
		if err != nil {
			return fmt.Errorf("Session.End: %w", err)
		}
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
