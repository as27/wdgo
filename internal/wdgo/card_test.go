package wdgo

import (
	"fmt"
	"testing"
)

func TestCardEventName(t *testing.T) {
	c := &Card{}
	c.Event(Cmd{"Name", "New Name"})
	if c.Name != "New Name" {
		t.Errorf("expect: New Name\ngot:%s", c.Name)
	}
}

func TestCardEventFields(t *testing.T) {
	c := &Card{}
	c.Event(Cmd{"Name", "New Name"})
	c.Event(Cmd{"Description", "Hello you"})
	c.Event(Cmd{"SupportID", "1223"})

	if c.Description != "Hello you" {
		t.Errorf("expect description: Hello you\ngot:%s", c.Name)
	}
	if c.SupportID != "1223" {
		t.Errorf("expect SupportID: 1223\ngot:%s", c.Name)
	}
}

func TestCardEventMoveTo(t *testing.T) {
	b := Board{}
	b.ID = "1"
	b.Name = "Testboard"
	b.Event(Cmd{"AddStage", "s0"})
	b.Event(Cmd{"AddStage", "s1"})
	e, _ := b.Find("s1")
	e.Event(Cmd{"AddCard", "c0"})
	e.Event(Cmd{"AddCard", "c1"})
	e.Event(Cmd{"AddCard", "c2"})
	e.Event(Cmd{"AddCard", "c3"})
	e.Event(Cmd{"AddCard", "c4"})
	e, _ = b.Find("c3")
	e.Event(Cmd{"MoveTo", "1"})
	s1 := b.Stages[1]
	if s1.Cards[0].ID != "c0" ||
		s1.Cards[1].ID != "c3" ||
		s1.Cards[2].ID != "c1" ||
		s1.Cards[3].ID != "c2" ||
		s1.Cards[4].ID != "c4" {
		t.Error("CardMoveTo error")
		for i, v := range s1.Cards {
			fmt.Println(i, v.ID)
		}
	}
}
