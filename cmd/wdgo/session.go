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
	case tcell.KeyTAB:

	}
	return event
}

func (a *app) renderSession() {
	const (
		startFormat = "02.01.06 15:04"
		endFormat   = "15:04"
	)
	a.card.sessions.Clear()
	activeBoard := &a.boards[a.activeBoard]
	activeStage := activeBoard.board.Stages[activeBoard.activeStage]
	activeCard := activeStage.Cards[activeBoard.activeCard]
	var sum time.Duration
	for i := range activeCard.Sessions {
		s := activeCard.Sessions[len(activeCard.Sessions)-i-1]
		a.card.sessions.SetCell(i, 0,
			tview.NewTableCell(s.Start.Format(startFormat)))
		if (s.End != time.Time{}) {
			a.card.sessions.SetCell(i, 1,
				tview.NewTableCell(s.End.Format(endFormat)))
			a.card.sessions.SetCell(i, 2,
				tview.NewTableCell(s.End.Sub(s.Start).String()))
			sum += s.End.Sub(s.Start)
		}
		a.card.sessions.SetCell(i, 3, tview.NewTableCell(s.Note))
	}
	a.card.sessions.SetCell(len(activeCard.Sessions), 2,
		tview.NewTableCell("sum:"))
	a.card.sessions.SetCell(len(activeCard.Sessions), 2,
		tview.NewTableCell(sum.String()))
	a.card.sessions.SetBorder(true)
	a.card.sessions.SetSelectable(true, false)
	a.card.sessions.SetSelectedFunc(func(row, col int) {
		a.renderSessionForm(row)
		a.root.SetFocus(a.card.sessionForm)
	})
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

func (a *app) renderSessionForm(index int) {
	const dateFormat = "2006-01-02 15:04"
	activeBoard := &a.boards[a.activeBoard]
	activeStage := activeBoard.board.Stages[activeBoard.activeStage]
	activeCard := activeStage.Cards[activeBoard.activeCard]
	if activeBoard.cardSelected == nil {
		activeBoard.cardSelected = activeCard
		//a.card.sessionForm.Clear(true)
		//return
	}
	sessions := len(activeCard.Sessions)
	i := sessions - index - 1
	if i < 0 {
		i = 0
	}
	session := activeCard.Sessions[i]
	startEdited, endEdited := false, false
	edited := wdgo.Session{}
	edited.Note = session.Note
	a.card.sessionForm.Clear(true)
	a.card.sessionForm.AddInputField("Start", session.Start.Format(dateFormat), 17,
		func(textToCheck string, lastChar rune) bool {
			return true
		},
		func(text string) {
			startEdited = true
			var err error
			edited.Start, err = time.Parse(dateFormat, text)
			if err != nil {
				startEdited = false
			}
		})
	a.card.sessionForm.AddInputField("End", session.End.Format(dateFormat), 17,
		func(textToCheck string, lastChar rune) bool {
			return true
		},
		func(text string) {
			endEdited = true
			var err error
			edited.End, err = time.Parse(dateFormat, text)
			if err != nil {
				endEdited = false
			}
		})
	a.card.sessionForm.AddInputField("Note", session.Note, 27,
		func(textToCheck string, lastChar rune) bool {
			return true
		},
		func(text string) {
			edited.Note = text
		})
	a.card.sessionForm.AddButton("Save", func() {
		if startEdited {
			activeBoard.aggregator.NewEvent(session.ID(), "Start",
				edited.Start.Format(wdgo.TimeFormat))
		}
		if endEdited {
			activeBoard.aggregator.NewEvent(session.ID(), "End",
				edited.End.Format(wdgo.TimeFormat))
		}
		if session.Note != edited.Note {
			activeBoard.aggregator.NewEvent(session.ID(), "Note", edited.Note)
		}
		activeBoard.aggregator.State()
		a.renderCard("edit")
		a.renderSession()
		a.card.sessionForm.Clear(true)
		a.root.SetFocus(a.card.form)
	})
	a.card.sessionForm.AddButton("Cancel", func() {
		a.card.sessionForm.Clear(true)
		a.root.SetFocus(a.card.form)
	})

}
