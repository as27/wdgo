package wdgo

import (
	"fmt"
	"time"
)

// Comment can be added to a card.
type Comment struct {
	EventSource
	Card *Card
	Text string
	Time time.Time
}

// Event implements the eventer-interface for the comments
// Following actions are implemented:
//  * Text
func (c *Comment) Event(cmd Cmd) error {
	switch cmd.Action {
	case "Text":
		c.Text = cmd.Value
	default:
		return fmt.Errorf("%w: %s", ErrNoSuchAction, cmd.Action)
	}
	return nil
}

// Find searches for an id inside the comment
func (c *Comment) Find(id string) (Eventer, error) {
	if c.ID == id {
		return c, nil
	}
	return nil, ErrIDNotFound
}