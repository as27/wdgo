package wdgo

import (
	"errors"
	"fmt"
	"testing"
)

func TestAddCard(t *testing.T) {
	b := Board{}
	b.ID = "1"
	b.Name = "Testboard"
	b.Event(Cmd{"AddStage", "s0"})
	b.Event(Cmd{"AddStage", "s1"})
	err := b.Stages[0].Event(Cmd{"AddCard", "c0"})
	expectNoError(t, err)
	err = b.Stages[0].Event(Cmd{"AddCard", "c1"})
	expectNoError(t, err)
	err = b.Stages[0].Event(Cmd{"AddCard", "c2"})
	expectNoError(t, err)
	if b.Stages[0].Cards[2].ID != "c2" {
		t.Errorf("expect ID=c2, got: %s", b.Stages[0].Cards[2].ID)
	}
}

func TestFind(t *testing.T) {
	b := Board{}
	b.ID = "1"
	b.Name = "Testboard"
	b.Event(Cmd{"AddStage", "s0"})
	b.Event(Cmd{"AddStage", "s1"})
	e, _ := b.Find("s1")
	e.Event(Cmd{"AddCard", "c0"})
	e.Event(Cmd{"AddCard", "c1"})
	ec1, err := b.Find("c1")
	expectNoError(t, err)
	c1, ok := ec1.(*Card)
	if !ok {
		t.Error("got wrong type")
		return
	}
	if c1.ID != "c1" {
		t.Errorf("expect c1\ngot: %s", c1.ID)
	}
	_, err = b.Find("cx")
	if !errors.Is(err, ErrIDNotFound) {
		t.Errorf("expect IDNotFound\ngot:%s", err)
	}
}

func TestStageMoveTo(t *testing.T) {
	b := Board{}
	b.ID = "1"
	b.Name = "Testboard"
	b.Event(Cmd{"AddStage", "s0"})
	b.Event(Cmd{"AddStage", "s1"})
	b.Event(Cmd{"AddStage", "s2"})
	b.Event(Cmd{"AddStage", "s3"})
	b.Event(Cmd{"AddStage", "s4"})
	e, _ := b.Find("s4")
	err := e.Event(Cmd{"MoveTo", "0"})
	expectNoError(t, err)
	if b.Stages[0].ID != "s4" && b.Stages[3].ID != "s2" {
		t.Errorf("MoveTo does not work\nGot: %#v", *b.Stages[0])
		for i, v := range b.Stages {
			fmt.Println(i, v.ID)
		}
	}
}
func XXXTest(t *testing.T) {
	b := Board{}
	b.ID = "1"
	b.Name = "Testboard"
	b.Event(Cmd{"AddStage", "s0"})
	b.Event(Cmd{"AddStage", "s1"})
	b.Event(Cmd{"AddStage", "s2"})
	b.Event(Cmd{"AddStage", "s3"})
	b.Event(Cmd{"AddStage", "s4"})
	//t.Errorf("%v", b)
}
