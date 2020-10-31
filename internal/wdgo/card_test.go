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
		t.Errorf("expect Support.id 1223\ngot:%s", c.Name)
	}
}

func TestCardEventMoveTo(t *testing.T) {
	b := Board{}
	b.id = "1"
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
	if s1.Cards[0].id != "c0" ||
		s1.Cards[1].id != "c3" ||
		s1.Cards[2].id != "c1" ||
		s1.Cards[3].id != "c2" ||
		s1.Cards[4].id != "c4" {
		t.Error("CardMoveTo error")
		for i, v := range s1.Cards {
			fmt.Println(i, v.id)
		}
	}
}

func TestAddSession(t *testing.T) {
	b := Board{}
	b.id = "1"
	b.Name = "Testboard"
	b.Event(Cmd{"AddStage", "s0"})
	b.Event(Cmd{"AddStage", "s1"})
	e, _ := b.Find("s1")
	e.Event(Cmd{"AddCard", "c0"})
	e.Event(Cmd{"AddCard", "c1"})
	e, _ = b.Find("c1")
	e.Event(Cmd{"AddSession", "ses0"})
	e.Event(Cmd{"AddSession", "ses1"})
	e.Event(Cmd{"AddSession", "ses2"})
	c1 := b.Stages[1].Cards[1]
	if len(c1.Sessions) != 3 ||
		c1.Sessions[0].id != "ses0" ||
		c1.Sessions[1].id != "ses1" ||
		c1.Sessions[2].id != "ses2" {
		t.Error("sessions are not added")
	}
}

func TestAddComment(t *testing.T) {
	b := Board{}
	b.id = "1"
	b.Name = "Testboard"
	b.Event(Cmd{"AddStage", "s0"})
	b.Event(Cmd{"AddStage", "s1"})
	e, _ := b.Find("s1")
	e.Event(Cmd{"AddCard", "c0"})
	e.Event(Cmd{"AddCard", "c1"})
	e, _ = b.Find("c1")
	e.Event(Cmd{"AddComment", "com0"})
	e.Event(Cmd{"AddComment", "com1"})
	c1 := b.Stages[1].Cards[1]
	if len(c1.Comments) != 2 ||
		c1.Comments[1].id != "com1" {
		t.Error("Comments not added correctly")
	}
}

func TestCardFind(t *testing.T) {
	card := &Card{}
	card.id = "card0"
	card.Event(Cmd{"AddComment", "com0"})
	card.Event(Cmd{"AddComment", "com1"})
	card.Event(Cmd{"AddSession", "ses0"})
	card.Event(Cmd{"AddSession", "ses1"})
	checkIDs := []string{"card0", "com0", "com1", "ses0", "ses1"}
	for _, id := range checkIDs {
		t.Run(id, func(t *testing.T) {
			e, err := card.Find(id)
			if err != nil {
				t.Errorf("no error expected. got: %s", err)
			}
			gotID, _ := EventerID(e)
			if gotID != id {
				t.Errorf("expect .id %s\ngot: %s", id, gotID)
			}
		})
	}
}
