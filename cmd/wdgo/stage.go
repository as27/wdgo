package main

import (
	"github.com/as27/wdgo/internal/wdgo"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
)

func (a *app) stageEvents(event *tcell.EventKey) *tcell.EventKey {
	//activeBoard := &a.boards[a.activeBoard]

	switch event.Key() {
	case tcell.KeyEsc:
		a.pages.SwitchToPage("board")
		a.root.SetFocus(a.home)
	}
	return event
}

func (a *app) renderEditStage(mode string) error {
	activeBoard := &a.boards[a.activeBoard]
	activeStage := activeBoard.board.Stages[activeBoard.activeStage]
	if mode == "add" {
		activeStage = &wdgo.Stage{}
	}
	a.stage.Clear(true)
	edited := *activeStage

	a.stage.AddInputField("Name", activeStage.Name, 20, nil,
		func(text string) { edited.Name = text })

	a.stage.AddButton("Save", func() {
		if mode == "add" {
			id := uuid.New().String()
			activeBoard.aggregator.NewEvent(activeBoard.board.ID(), "AddStage", id)
			activeBoard.aggregator.NewEvent(id, "Name", edited.Name)
		} else {
			activeBoard.aggregator.NewEvent(activeStage.ID(), "Name", edited.Name)
		}
		activeBoard.aggregator.State()
		a.renderBoard()
	})
	a.stage.AddButton("Cancel", func() {
		a.renderBoard()
	})
	a.stage.SetTitle("stage properties").SetBorder(true)
	a.pages.AddAndSwitchToPage("stage", a.stage, true)
	a.root.SetFocus(a.stage)
	return nil
}
