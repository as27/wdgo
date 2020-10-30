package wdgo

import (
	"testing"
	"time"
)

func TestEventSourceAction(t *testing.T) {
	now := time.Now()
	later := now.Add(time.Second)
	es := EventSource{ID: "1", Name: "some name"}
	es.Action(now)
	if !es.Created.Equal(now) {
		t.Errorf("1: Created\nexpect:%s\ngot:%s", now, es.Created)
	}
	if !es.Modified.Equal(now) {
		t.Errorf("1: Modified\nexpect:%s\ngot:%s", now, es.Modified)
	}
	es.Action(later)
	if !es.Created.Equal(now) {
		t.Errorf("2: Created\nexpect:%s\ngot:%s", now, es.Created)
	}
	if !es.Modified.Equal(later) {
		t.Errorf("2: Modified\nexpect:%s\ngot:%s", now, es.Modified)
	}
}

func TestEventerID(t *testing.T) {
	exp := func(t *testing.T, b Board, id string) {
		t.Run(id, func(t *testing.T) {
			e, err := b.Find(id)
			if err != nil {
				t.Errorf("Find(%s):no error expected\ngot: %s", id, err)
				t.FailNow()
			}
			gotID, err := EventerID(e)
			if err != nil {
				t.Errorf("no error expected\ngot: %s", err)
				t.FailNow()
			}
			if id != gotID {
				t.Errorf("expect: %s got: %s", id, gotID)
			}
		})
	}
	b := Board{}
	b.ID = "board0"
	b.Name = "Testboard"
	b.Event(Cmd{"AddStage", "s0"})
	b.Stages[0].Event(Cmd{"AddCard", "c0"})
	e, _ := b.Find("c0")
	e.Event(Cmd{"AddSession", "ss0"})
	e.Event(Cmd{"AddComment", "cc0"})
	exp(t, b, "board0")
	exp(t, b, "s0")
	exp(t, b, "c0")
	exp(t, b, "ss0")
	exp(t, b, "cc0")
}
