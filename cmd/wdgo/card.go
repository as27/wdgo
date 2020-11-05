package main

import (
	"github.com/gdamore/tcell/v2"
)

func (a *app) cardEvents(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlQ:
		a.renderBoard()
	}
	return event
}

func (a *app) renderCard() error {
	activeBoard := &a.boards[a.activeBoard]
	activeStage := activeBoard.board.Stages[activeBoard.activeStage]
	activeCard := activeStage.Cards[activeBoard.activeCard]
	a.card.Clear(true)
	edited := *activeCard
	stages := []string{}
	for _, s := range activeBoard.board.Stages {
		stages = append(stages, s.Name)
	}
	a.card.AddInputField("Name", activeCard.Name, 20, nil,
		func(text string) { edited.Name = text })
	a.card.AddInputField("Description", activeCard.Description, 20, nil,
		func(text string) { edited.Description = text })
	a.card.AddInputField("Support ID", activeCard.SupportID, 20, nil,
		func(text string) { edited.SupportID = text })
	a.card.AddInputField("Customer", activeCard.Customer, 20, nil,
		func(text string) { edited.Customer = text })
	a.card.AddDropDown("Stage", stages, activeBoard.activeStage, nil)
	a.card.AddButton("Save", func() {
		*activeCard = edited
		a.renderBoard()
	})
	a.card.AddButton("Cancel", func() {
		a.renderBoard()
	})

	a.pages.AddAndSwitchToPage("card", a.card, true)
	a.root.SetFocus(a.card)
	return nil
}
