package main

import (
	"github.com/as27/wdgo/internal/wdgo"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
)

func (a *app) cardEvents(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlQ:
		a.renderBoard()
	}
	return event
}

func (a *app) renderCard(mode string) error {
	activeBoard := &a.boards[a.activeBoard]
	activeStage := activeBoard.board.Stages[activeBoard.activeStage]
	activeCard := &wdgo.Card{}
	if mode == "edit" {
		activeCard = activeStage.Cards[activeBoard.activeCard]
	}
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
		}
		activeBoard.aggregator.State()
		a.renderBoard()
	})
	a.card.AddButton("Cancel", func() {
		a.renderBoard()
	})

	a.pages.AddAndSwitchToPage("card", a.card, true)
	a.root.SetFocus(a.card)
	return nil
}
