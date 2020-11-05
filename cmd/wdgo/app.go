package main

import (
	"fmt"
	"log"

	"github.com/as27/wdgo/internal/wdgo"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type app struct {
	root        *tview.Application
	pages       *tview.Pages
	home        *tview.List
	activeBoard int
	boards      []board
}

type board struct {
	view        *tview.Flex
	board       *wdgo.Board
	stageviews  []*tview.Flex
	cardviews   [][]*tview.TextView
	activeStage int
	activeCard  int
}

func newApp() *app {
	a := app{
		root:  tview.NewApplication(),
		pages: tview.NewPages(),
		home:  tview.NewList(),
	}
	a.initBoards()
	return &a
}

func (a *app) run() error {
	err := a.layout()
	if err != nil {
		return fmt.Errorf("app.run().layout(): %w", err)
	}
	return a.root.Run()
}

func (a *app) layout() error {
	a.renderHome()
	a.pages.AddPage("board", tview.NewBox(), true, false)
	a.pages.AddPage("home", a.home, true, true)
	a.root.SetRoot(a.pages, true).EnableMouse(true)
	a.root.SetInputCapture(a.inputCaptures)
	return nil
}

func (a *app) inputCaptures(event *tcell.EventKey) *tcell.EventKey {
	pageName, _ := a.pages.GetFrontPage()
	log.Println(pageName)
	switch pageName {
	case "home":
		a.homeEvents(event)
	case "board":
		a.boardEvents(event)
	case "card":
		a.cardEvents(event)
	}
	// global key bindings
	if event.Key() == tcell.KeyEsc {
		a.root.Stop()
	}
	return event
}

func (a *app) cardEvents(event *tcell.EventKey) *tcell.EventKey {
	return event
}

func (a *app) setCard() error {
	//	a.view.card = tview.NewForm()
	//	a.view.card.SetTitle("Card")
	return nil
}
func (a *app) setFoo() error { return nil }
