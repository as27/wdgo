package wdgo

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	ErrNoSuchAction = errors.New("no such a action")
)

type Eventer interface {
	Event(Cmd) error
}

type Cmd struct {
	Action string
	Value  string
}

// Board defines a Kanban-board
type Board struct {
	ID     string
	Name   string
	Stages []Stage
}

func (b *Board) Event(cmd Cmd) error {
	switch cmd.Action {
	case "Name":
		b.Name = cmd.Value
	case "AddStage":
		s := Stage{
			ID:    cmd.Value,
			Board: b,
		}
		b.Stages = append(b.Stages, s)
	default:
		return fmt.Errorf("%w: %s", ErrNoSuchAction, cmd.Action)
	}
	return nil
}

func (b *Board) Find(id string) (Eventer, error) {
	return nil, nil
}

// Stage is always part of a board
type Stage struct {
	ID    string
	Name  string
	Board *Board
	Cards []Card
}

func (s *Stage) Event(cmd Cmd) error {
	switch cmd.Action {
	case "Name":
		s.Name = cmd.Value
	case "AddCard":
		c := Card{
			ID:    cmd.Value,
			Stage: s,
		}
		s.Cards = append(s.Cards, c)
	case "MoveTo":
		pos, err := strconv.Atoi(cmd.Value)
		if err != nil {
			return fmt.Errorf("Stage.Event(MoveTo:%s):%w", cmd.Value, err)
		}
		for i, v := range s.Board.Stages {
			if i == pos {
				s.Board.Stages = append(s.Board.Stages, *s)
			}
			if v.ID != s.ID {
				s.Board.Stages = append(s.Board.Stages, v)
			}
		}
	default:
		return fmt.Errorf("%w: %s", ErrNoSuchAction, cmd.Action)
	}
	return nil
}

func (s Stage) Find(id string) (Eventer, error) {
	return nil, nil
}

// Card is part of a stage.
type Card struct {
	ID          string
	Stage       *Stage
	Title       string
	Description string
	SupportID   string
	Comments    []Comment
}

func (c Card) Find(id string) (Eventer, error) {
	return nil, nil
}

// Session is used to track the time, when working on a card
type Session struct {
	ID    string
	Card  *Card
	Title string
	Start time.Time
	End   time.Time
}

// Comment can be added to a card.
type Comment struct {
	ID string

	Card *Card
	Text string
	Time time.Time
}
