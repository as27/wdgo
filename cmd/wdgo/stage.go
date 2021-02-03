package main

import (
	"fmt"
	"strconv"

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

	var stages []string
	for i, s := range activeBoard.board.Stages {
		stages = append(stages, fmt.Sprintf("%0d (%s)", i, s.Name))
	}
	stageIndex := activeBoard.activeStage

	a.stage.AddInputField("Name", activeStage.Name, 20, nil,
		func(text string) { edited.Name = text })

	a.stage.AddDropDown("Position", stages, stageIndex, func(option string, optionIndex int) {
		stageIndex = optionIndex
	})

	a.stage.AddButton("Save", func() {
		if mode == "add" {
			id := uuid.New().String()
			activeBoard.aggregator.NewEvent(activeBoard.board.ID(), "AddStage", id)
			activeBoard.aggregator.NewEvent(id, "Name", edited.Name)
		} else {
			if edited.Name != activeStage.Name {
				activeBoard.aggregator.NewEvent(activeStage.ID(), "Name", edited.Name)
			}
		}
		if stageIndex != activeBoard.activeStage {
			activeBoard.aggregator.NewEvent(
				activeStage.ID(),
				"MoveTo", strconv.Itoa(stageIndex))
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
