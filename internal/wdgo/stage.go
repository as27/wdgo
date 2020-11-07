package wdgo

import (
	"errors"
	"fmt"
	"strconv"
)

// Stage is always part of a board
type Stage struct {
	EventSource
	Board *Board
	Cards []*Card
}

// NewStage creates a stage with an inmutable id
func NewStage(id string) *Stage {
	s := Stage{}
	s.id = id
	return &s
}

// Event implements the Eventer-interface.
// Following actions are implemented:
//  * Name
//  * AddCard
//  * MoveTo
func (s *Stage) Event(cmd Cmd) error {
	switch cmd.Action {
	case "Name":
		s.Name = cmd.Value
	case "AddCard":
		c := NewCard(cmd.Value)
		c.Stage = s
		s.Cards = append(s.Cards, c)
	case "MoveTo":
		pos, err := strconv.Atoi(cmd.Value)
		if err != nil {
			return fmt.Errorf("Stage.Event(MoveTo:%s):%w", cmd.Value, err)
		}
		stages := []*Stage{}
		for i, v := range s.Board.Stages {
			if i == pos {
				stages = append(stages, s)
			}
			if v.id != s.id {
				stages = append(stages, v)
			}
		}
		s.Board.Stages = stages
	default:
		return fmt.Errorf("%w: %s", ErrNoSuchAction, cmd.Action)
	}
	return nil
}

// Find searchs for an id inside the stage. That also contains the
// cards, comments, sessions
func (s *Stage) Find(id string) (Eventer, error) {
	if s.id == id {
		return s, nil
	}
	for _, c := range s.Cards {
		e, err := c.Find(id)
		if err == nil {
			return e, nil
		}
		if !errors.Is(err, ErrIDNotFound) {
			return nil, fmt.Errorf("Board.Find(%s):%w", id, err)
		}
	}
	return nil, ErrIDNotFound
}

// Pos returns the position inside the board
func (s *Stage) Pos() int {
	for i, ss := range s.Board.Stages {
		if s.ID() == ss.ID() {
			return i
		}
	}
	return -1
}
