package main

import (
	"github.com/as27/wdgo/internal/wdgo"
	"github.com/google/uuid"
)

// This file is just for development purpose until
// the elements can be added over the gui

func (a *app) initBoardsTest() {
	b := wdgo.NewBoard(uuid.New().String())
	b.Name = "Board 1"
	id := uuid.New().String()
	b.Event(wdgo.Cmd{"AddStage", id})
	s, _ := b.Find(id)
	s.Event(wdgo.Cmd{"Name", "Backlog"})
	s.Event(wdgo.Cmd{"AddCard", "c1"})
	c, _ := b.Find("c1")
	c.Event(wdgo.Cmd{"Name", "Todo 1"})
	s.Event(wdgo.Cmd{"AddCard", "c2"})
	c, _ = b.Find("c2")
	c.Event(wdgo.Cmd{"Name", "Todo 2"})
	s.Event(wdgo.Cmd{"AddCard", "c3"})
	c, _ = b.Find("c3")
	c.Event(wdgo.Cmd{"Name", "Todo 3"})

	id = uuid.New().String()
	b.Event(wdgo.Cmd{"AddStage", id})
	s, _ = b.Find(id)
	s.Event(wdgo.Cmd{"Name", "Doing"})
	s.Event(wdgo.Cmd{"AddCard", "c11"})
	c, _ = b.Find("c11")
	c.Event(wdgo.Cmd{"Name", "Todo 1"})

	id = uuid.New().String()
	b.Event(wdgo.Cmd{"AddStage", id})
	s, _ = b.Find(id)
	s.Event(wdgo.Cmd{"Name", "Done"})
	a.addBoard(b)
	b = wdgo.NewBoard(uuid.New().String())
	b.Name = "Board 2"
	a.addBoard(b)
}
