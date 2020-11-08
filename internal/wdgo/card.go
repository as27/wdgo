package wdgo

import (
	"errors"
	"fmt"
	"strconv"
)

// Card is part of a stage.
type Card struct {
	EventSource
	Stage       *Stage
	Description string
	SupportID   string
	Customer    string
	Comments    []*Comment
	Sessions    []*Session
}

// NewCard creates a card with a given id
func NewCard(id string) *Card {
	c := Card{}
	c.id = id
	return &c
}

// Event implements the Eventer-interface.
// Following actions are implemented:
//  * Name
//  * Description
//  * SupportID
//  * Customer
//  * MoveTo
//  * AddSession
//  * AddComment
func (c *Card) Event(cmd Cmd) error {
	switch cmd.Action {
	case "Name":
		c.Name = cmd.Value
	case "Description":
		c.Description = cmd.Value
	case "SupportID":
		c.SupportID = cmd.Value
	case "Customer":
		c.Customer = cmd.Value
	case "MoveTo":
		pos, err := strconv.Atoi(cmd.Value)
		if err != nil {
			return fmt.Errorf("Card.Event(MoveTo:%s):%w", cmd.Value, err)
		}
		oldPos := c.Pos()
		cards := []*Card{}
		// remove this card from stack
		c.Stage.Cards = append(
			c.Stage.Cards[:oldPos], c.Stage.Cards[oldPos+1:]...)
		// Add the card at new position
		for i, v := range c.Stage.Cards {
			if i == pos {
				cards = append(cards, c)
			}
			if v.id != c.id {
				cards = append(cards, v)
			}
		}
		// check if card moved to last position
		if pos == len(c.Stage.Cards) {
			cards = append(cards, c)
		}
		c.Stage.Cards = cards
	case "MoveToStage":
		e, err := c.Stage.Board.Find(cmd.Value)
		if err != nil {
			return fmt.Errorf("MoveToStage: %s not found: %w",
				cmd.Value, err)
		}
		stage, ok := e.(*Stage)
		if !ok {
			return fmt.Errorf("id: %s is not a stage", cmd.Value)
		}
		stage.Cards = append(stage.Cards, c)
		cards := []*Card{}
		for _, cc := range c.Stage.Cards {
			if c.ID() == cc.ID() {
				continue
			}
			cards = append(cards, cc)
		}
		c.Stage.Cards = cards // set the cards of the old stage
		c.Stage = stage       // internal point to the new stage
	case "AddSession":
		s := NewSession(cmd.Value)
		s.Card = c
		c.Sessions = append(c.Sessions, s)
	case "AddComment":
		cc := NewComment(cmd.Value)
		cc.Card = c
		c.Comments = append(c.Comments, cc)
	default:
		return fmt.Errorf("%w: %s", ErrNoSuchAction, cmd.Action)
	}
	return nil
}

// Find searches an id inside this card, sessions and comments.
func (c *Card) Find(id string) (Eventer, error) {
	if c.id == id {
		return c, nil
	}
	for _, s := range c.Sessions {
		e, err := s.Find(id)
		if err == nil {
			return e, nil
		}
		if !errors.Is(err, ErrIDNotFound) {
			return nil, fmt.Errorf("Card.Find(%s):%w", id, err)
		}
	}
	for _, cc := range c.Comments {
		e, err := cc.Find(id)
		if err == nil {
			return e, nil
		}
		if !errors.Is(err, ErrIDNotFound) {
			return nil, fmt.Errorf("Card.Find(%s):%w", id, err)
		}
	}
	return nil, ErrIDNotFound
}

// Pos returns the position inside the stage
func (c *Card) Pos() int {
	for i, cc := range c.Stage.Cards {
		if c.ID() == cc.ID() {
			return i
		}
	}
	return -1
}
