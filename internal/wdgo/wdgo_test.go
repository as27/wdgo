package wdgo

import (
	"testing"
	"time"
)

func TestEventSourceAction(t *testing.T) {
	now := time.Now()
	later := now.Add(time.Second)
	es := EventSource{id: "1", Name: "some name"}
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
