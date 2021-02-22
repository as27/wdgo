package wdgo

import "testing"

func TestSessionEvents(t *testing.T) {
	b := Board{}
	b.id = "1"
	b.Name = "Testboard"
	b.Event(Cmd{"AddStage", "s0"})
	b.Event(Cmd{"AddStage", "s1"})
	e, _ := b.Find("s1")
	e.Event(Cmd{"AddCard", "c0"})
	e.Event(Cmd{"AddCard", "c1"})
	e, _ = b.Find("c1")
	e.Event(Cmd{"AddSession", "ss0"})
	err := e.Event(Cmd{"AddSession", "ss1"})
	if err != nil {
		t.Errorf("AddSession:no error expected.\ngot: %s", err)
	}
	e, err = b.Find("ss1")
	if err != nil {
		t.Errorf("Find(ss1)\nno error expected.\ngot: %s", err)
	}
	expectNote := "A 123 Note"
	err = e.Event(Cmd{"Note", expectNote})
	expectNoError(t, err)
	gotNote := b.Stages[1].Cards[1].Sessions[1].Note
	if gotNote != expectNote {
		t.Errorf("ss1 Note should: %s; got: %s", expectNote, gotNote)
	}
}
