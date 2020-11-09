package main

import (
	"time"

	"github.com/as27/wdgo/internal/wdgo"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func (a *app) sessionEvents(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlJ:
		// start session
	case tcell.KeyCtrlK:
		// stop running session
	}
	return event
}

func (a *app) renderSession() {
	a.card.sessions.Clear()
	activeBoard := &a.boards[a.activeBoard]
	activeStage := activeBoard.board.Stages[activeBoard.activeStage]
	activeCard := activeStage.Cards[activeBoard.activeCard]
	for i := range activeCard.Sessions {
		s := activeCard.Sessions[len(activeCard.Sessions)-i-1]
		a.card.sessions.SetCell(i, 0,
			tview.NewTableCell(s.Start.Format(wdgo.TimeFormat)))
		if (s.End != time.Time{}) {
			a.card.sessions.SetCell(i, 1,
				tview.NewTableCell(s.End.Format(wdgo.TimeFormat)))
			a.card.sessions.SetCell(i, 2,
				tview.NewTableCell(s.End.Sub(s.Start).String()))
		}
	}
	a.card.sessions.SetBorder(true)
	a.card.sessions.SetBorders(true)
}
