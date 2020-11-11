package main

import (
	"time"

	"github.com/as27/wdgo/internal/wdgo"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

func (a *app) sessionEvents(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlJ:
		// start/stop session
		a.sessionStartStop()
	}
	return event
}

func (a *app) renderSession() {
	a.card.sessions.Clear()
	activeBoard := &a.boards[a.activeBoard]
	activeStage := activeBoard.board.Stages[activeBoard.activeStage]
	activeCard := activeStage.Cards[activeBoard.activeCard]
	var sum time.Duration
	for i := range activeCard.Sessions {
		s := activeCard.Sessions[len(activeCard.Sessions)-i-1]
		a.card.sessions.SetCell(i, 0,
			tview.NewTableCell(s.Start.Format(wdgo.TimeFormat)))
		if (s.End != time.Time{}) {
			a.card.sessions.SetCell(i, 1,
				tview.NewTableCell(s.End.Format(wdgo.TimeFormat)))
			a.card.sessions.SetCell(i, 2,
				tview.NewTableCell(s.End.Sub(s.Start).String()))
			sum += s.End.Sub(s.Start)
		}
	}
	a.card.sessions.SetCell(len(activeCard.Sessions), 2,
		tview.NewTableCell("sum:"))
	a.card.sessions.SetCell(len(activeCard.Sessions), 2,
		tview.NewTableCell(sum.String()))
	a.card.sessions.SetBorder(true)
	//a.card.sessions.SetBorders(true)
	//a.card.sessions.scro
	a.card.sessions.ScrollToBeginning()

}

func (a *app) sessionStartStop() {
	activeBoard := &a.boards[a.activeBoard]
	activeStage := activeBoard.board.Stages[activeBoard.activeStage]
	activeCard := activeStage.Cards[activeBoard.activeCard]
	now := time.Now().Format(wdgo.TimeFormat)
	if (len(activeCard.Sessions) == 0 ||
		activeCard.Sessions[len(activeCard.Sessions)-1].End != time.Time{}) {
		// create new session
		id := uuid.New().String()
		activeBoard.aggregator.NewEvent(activeCard.ID(), "AddSession", id)
		activeBoard.aggregator.NewEvent(id, "Start", now)
	} else {
		id := activeCard.Sessions[len(activeCard.Sessions)-1].ID()
		activeBoard.aggregator.NewEvent(id, "End", now)
	}
	activeBoard.aggregator.State()
	a.renderCard("edit")
}
