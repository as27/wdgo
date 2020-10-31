package wdgo

import "testing"

func TestCommentEvents(t *testing.T) {
	c := &Comment{}
	c.id = "cc0"
	txt := "some text"
	c.Event(Cmd{"Text", txt})
	if c.Text != txt {
		t.Errorf("Event.Text: want: %s got: %s", txt, c.Text)
	}
}
