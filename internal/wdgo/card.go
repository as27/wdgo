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
	Comments    []*Comment
	Sessions    []*Session
}

// Event implements the Eventer-interface.
// Following actions are implemented:
//  * Name
//  * Description
//  * SupportID
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
	case "MoveTo":
		pos, err := strconv.Atoi(cmd.Value)
		if err != nil {
			return fmt.Errorf("Card.Event(MoveTo:%s):%w", cmd.Value, err)
		}
		cards := []*Card{}
		for i, v := range c.Stage.Cards {
			if i == pos {
				cards = append(cards, c)
			}
			if v.ID != c.ID {
				cards = append(cards, v)
			}
		}
		c.Stage.Cards = cards
	case "AddSession":
		s := &Session{}
		s.ID = cmd.Value
		s.Card = c
		c.Sessions = append(c.Sessions, s)
	case "AddComment":
		cc := &Comment{}
		cc.ID = cmd.Value
		cc.Card = c
		c.Comments = append(c.Comments, cc)
	default:
		return fmt.Errorf("%w: %s", ErrNoSuchAction, cmd.Action)
	}
	return nil
}

// Find searches an id inside this card, sessions and comments.
func (c *Card) Find(id string) (Eventer, error) {
	if c.ID == id {
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
