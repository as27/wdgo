package wdgo

import (
	"errors"
	"fmt"
)

// Board defines a Kanban-board
type Board struct {
	EventSource
	Stages []*Stage
}

// Event implements the eventer-interface for Board.
// It implements:
//   * Name
//   * AddStage
func (b *Board) Event(cmd Cmd) error {
	switch cmd.Action {
	case "Name":
		b.Name = cmd.Value
	case "AddStage":
		s := Stage{}
		s.ID = cmd.Value
		s.Board = b
		b.Stages = append(b.Stages, &s)
	default:
		return fmt.Errorf("%w: %s", ErrNoSuchAction, cmd.Action)
	}
	return nil
}

// Find an ID, which is part of the board. That method search
// also inside the stages, cards, comments and sessions.
func (b *Board) Find(id string) (Eventer, error) {
	if b.ID == id {
		return b, nil
	}
	for _, s := range b.Stages {
		e, err := s.Find(id)
		if err == nil {
			return e, nil
		}
		if !errors.Is(err, ErrIDNotFound) {
			return nil, fmt.Errorf("Board.Find(%s):%w", id, err)
		}
	}
	return nil, ErrIDNotFound
}
