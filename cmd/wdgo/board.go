package main

import (
	"io"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

func (a *app) boardEvents(event *tcell.EventKey) *tcell.EventKey {
	activeBoard := &a.boards[a.activeBoard]
	activeStage := activeBoard.board.Stages[activeBoard.activeStage]
	cardsInBoard := len(activeBoard.board.Stages[activeBoard.activeStage].Cards)

	switch event.Key() {
	case tcell.KeyBackspace:
		a.pages.SwitchToPage("home")
		a.root.SetFocus(a.home)
	case tcell.KeyTAB, tcell.KeyRight:
		activeBoard.activeStage++
		if activeBoard.activeStage >= len(activeBoard.board.Stages) {
			activeBoard.activeStage = 0
		}
		activeBoard.activeCard = 0
		a.renderBoard()
	case tcell.KeyBacktab, tcell.KeyLeft:
		activeBoard.activeStage--
		if activeBoard.activeStage < 0 {
			activeBoard.activeStage = len(activeBoard.board.Stages) - 1
		}
		activeBoard.activeCard = 0
		a.renderBoard()
	case tcell.KeyDown:
		if activeBoard.cardSelected != nil {
			pos := activeBoard.cardSelected.Pos()
			pos++
			if pos >= len(activeStage.Cards) {
				pos = len(activeStage.Cards) - 1
			}
			activeBoard.aggregator.NewEvent(
				activeBoard.cardSelected.ID(),
				"MoveTo", strconv.Itoa(pos))
			activeBoard.aggregator.State()
		}
		if activeBoard.activeCard < cardsInBoard-1 {
			activeBoard.activeCard++
		}
		a.renderBoard()
	case tcell.KeyUp:
		if activeBoard.cardSelected != nil {
			pos := activeBoard.cardSelected.Pos()
			if pos > 0 {
				pos--
				activeBoard.aggregator.NewEvent(
					activeBoard.cardSelected.ID(),
					"MoveTo", strconv.Itoa(pos))
				activeBoard.aggregator.State()
			}
		}
		if activeBoard.activeCard > 0 {
			activeBoard.activeCard--
		}
		a.renderBoard()
	case tcell.KeyCtrlE:
		if activeBoard.cardSelected != nil {
			a.renderCard("edit")
			break
		}
		// edit stage
		a.renderEditStage("edit")
	case tcell.KeyCtrlA:
		// add stage
		a.renderEditStage("add")
	case tcell.KeyCtrlN:
		// new card
		a.renderCard("add")
	case tcell.KeyEnter:
		if activeBoard.cardSelected != nil {
			activeBoard.cardSelected = nil
		} else {
			activeBoard.cardSelected = activeStage.Cards[activeBoard.activeCard]
		}
		a.renderBoard()
		//a.renderCard("edit")
	}
	return event
}

func (a *app) renderBoard() error {
	activeBoard := &a.boards[a.activeBoard]
	if activeBoard.view == nil {
		activeBoard.view = tview.NewFlex()
	}
	activeBoard.view.Clear()
	a.boards[a.activeBoard].stageviews = []*tview.Flex{}
	a.boards[a.activeBoard].cardviews = [][]*tview.TextView{}
	if len(activeBoard.board.Stages) == 0 {
		id := uuid.New().String()
		activeBoard.aggregator.NewEvent(
			activeBoard.board.ID(), "AddStage", id)
		activeBoard.aggregator.NewEvent(id, "Name", "New Stage")
		activeBoard.aggregator.State()
	}
	for snr, s := range activeBoard.board.Stages {
		stage := tview.NewFlex()
		stage.SetTitle(s.Name).SetBorder(true)
		active := false
		if snr == activeBoard.activeStage {
			stage.SetBorderColor(tcell.ColorGreenYellow)
			active = true
		}
		cards := []*tview.TextView{}
		cardbox := tview.NewFlex().SetDirection(tview.FlexRow)
		for cnr, c := range s.Cards {
			txt := tview.NewTextView()
			txt.SetBorder(true)
			if active && activeBoard.activeCard == cnr {
				txt.SetBorderColor(tcell.ColorGreenYellow)
				if activeBoard.cardSelected != nil {
					txt.SetBackgroundColor(tcell.ColorWhite)
					txt.SetTextColor(tcell.ColorBlack)
				}
			}
			io.WriteString(txt, c.Name)
			cards = append(cards, txt)
			cardbox.AddItem(txt, 3, 1, false)
		}
		stage.AddItem(cardbox, 0, 1, false)
		activeBoard.cardviews =
			append(activeBoard.cardviews, cards)
		a.boards[a.activeBoard].stageviews =
			append(a.boards[a.activeBoard].stageviews, stage)
	}
	for _, v := range a.boards[a.activeBoard].stageviews {
		activeBoard.view.AddItem(v, 0, 1, false)
	}
	a.pages.AddAndSwitchToPage("board", activeBoard.view, true)
	return nil
}
