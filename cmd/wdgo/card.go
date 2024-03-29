package main

import (
	"log"

	"github.com/as27/wdgo/internal/wdgo"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

// setActiveCard is used to set the focus to a card by using
// the id
func (a *app) setActiveCard(id string) (ok bool) {
	for ib, b := range a.boards {
		for is, s := range b.board.Stages {
			ic := 0
			for _, c := range s.Cards {

				if c.ID() == id {
					a.activeBoard = ib
					a.boards[ib].activeStage = is
					a.boards[ib].activeCard = ic
					return true
				}
				// just count valid cards, because archived
				// are not visible inside the app
				if !c.Archived {
					ic++
				}
			}
		}
	}
	return false
}

func (a *app) cardEvents(event *tcell.EventKey) *tcell.EventKey {
	activeBoard := &a.boards[a.activeBoard]
	//activeStage := activeBoard.board.Stages[activeBoard.activeStage]
	switch event.Key() {
	case tcell.KeyEsc:
		activeBoard.cardSelected = nil
		a.renderBoard()
		// store all events, when leaving a card view
		go func() {
			err := a.writeEvents()
			if err != nil {
				log.Printf("SwitchToHome: writeEvent() err: %s", err)
			}
		}()
	case tcell.KeyRight, tcell.KeyLeft:
		if a.card.form.HasFocus() {
			a.root.SetFocus(a.card.sessions)
		} else if a.card.sessions.HasFocus() {
			a.root.SetFocus(a.card.form)
		}
	}
	return event
}

func (a *app) renderCard(mode string) error {
	activeBoard := &a.boards[a.activeBoard]
	activeStage := activeBoard.board.Stages[activeBoard.activeStage]
	activeCard := &wdgo.Card{}
	if mode == "edit" {
		activeCard = activeStage.GetActiveCardNr(activeBoard.activeCard)
	}
	a.card.form.Clear(true)
	edited := *activeCard
	stages := []string{}
	for _, s := range activeBoard.board.Stages {
		stages = append(stages, s.Name)
	}
	stageIndex := activeBoard.activeStage
	if mode == "edit" {
		a.card.form.AddButton("Start/Stop Session", a.sessionStartStop)
	}
	a.card.form.AddInputField("Name", activeCard.Name, 20, nil,
		func(text string) { edited.Name = text })
	a.card.form.AddInputField("Description", activeCard.Description, 20, nil,
		func(text string) { edited.Description = text })
	a.card.form.AddInputField("Support ID", activeCard.SupportID, 20, nil,
		func(text string) { edited.SupportID = text })
	a.card.form.AddInputField("Customer", activeCard.Customer, 20, nil,
		func(text string) { edited.Customer = text })
	a.card.form.AddDropDown("Stage", stages, activeBoard.activeStage, func(option string, optionIndex int) {
		stageIndex = optionIndex
	})
	a.card.form.AddCheckbox("ToDo", activeCard.ToDo,
		func(checked bool) { edited.ToDo = checked })
	a.card.form.AddCheckbox("Archived", activeCard.Archived,
		func(checked bool) { edited.Archived = checked })
	a.card.form.AddButton("Save", func() {
		if mode == "add" {
			id := uuid.New().String()
			activeBoard.aggregator.NewEvent(activeStage.ID(), "AddCard", id)
			if edited.Name != "" {
				activeBoard.aggregator.NewEvent(id, "Name", edited.Name)
			}
			if edited.Description != "" {
				activeBoard.aggregator.NewEvent(id, "Description", edited.Description)
			}
			if edited.SupportID != "" {
				activeBoard.aggregator.NewEvent(id, "SupportID", edited.SupportID)
			}
			if edited.Customer != "" {
				activeBoard.aggregator.NewEvent(id, "Customer", edited.Customer)
			}
			if edited.ToDo {
				activeBoard.aggregator.NewEvent(id, "ToDo", "True")
			}
		} else {
			id := activeCard.ID()
			if edited.Name != activeCard.Name {
				activeBoard.aggregator.NewEvent(id, "Name", edited.Name)
			}
			if edited.Description != activeCard.Description {
				activeBoard.aggregator.NewEvent(id, "Description", edited.Description)
			}
			if edited.SupportID != activeCard.SupportID {
				activeBoard.aggregator.NewEvent(id, "SupportID", edited.SupportID)
			}
			if edited.Customer != activeCard.Customer {
				activeBoard.aggregator.NewEvent(id, "Customer", edited.Customer)
			}
			if edited.ToDo != activeCard.ToDo {
				val := "False"
				if edited.ToDo {
					val = "True"
				}
				activeBoard.aggregator.NewEvent(id, "ToDo", val)
			}
			if edited.Archived != activeCard.Archived {
				val := "False"
				if edited.Archived {
					val = "True"
				}
				activeBoard.aggregator.NewEvent(id, "Archived", val)
			}
		}
		if activeBoard.activeStage != stageIndex {
			activeBoard.aggregator.NewEvent(activeCard.ID(),
				"MoveToStage", activeBoard.board.Stages[stageIndex].ID())
			activeBoard.activeStage = stageIndex
			activeBoard.activeCard = 0
		}
		activeBoard.cardSelected = nil
		activeBoard.aggregator.State()
		a.renderBoard()
	})
	a.card.form.AddButton("Cancel", func() {
		activeBoard.cardSelected = nil
		a.renderBoard()
	})
	a.card.form.SetTitle("card properties").SetBorder(true)
	a.card.card.Clear()
	a.card.card.AddItem(a.card.form, 0, 1, true)
	if mode == "edit" {
		a.renderSession()
	} else {
		a.card.sessions.Clear()
	}
	a.card.sessionsFlex.Clear()
	a.card.sessionsFlex.SetDirection(tview.FlexRow)
	a.card.sessionsFlex.AddItem(a.card.sessions, 0, 1, false)
	a.card.sessionsFlex.AddItem(a.card.sessionForm, 10, 1, false)
	a.card.card.AddItem(a.card.sessionsFlex, 0, 1, false)
	//a.card.card.AddItem(a.card.sessions, 0, 1, false)
	a.pages.AddAndSwitchToPage("card", a.card.card, true)
	a.root.SetFocus(a.card.form)
	return nil
}
