package main

import (
	"fmt"
	"io"

	"github.com/as27/wdgo/internal/wdgo"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

type app struct {
	view struct {
		root  *tview.Application
		pages *tview.Pages
		home  *tview.List
		board *tview.Flex
		card  *tview.Form
	}
	boards        []*wdgo.Board
	selectedBoard *wdgo.Board
}

func (a *app) initBoards() {
	b := wdgo.NewBoard(uuid.New().String())
	b.Name = "Board 1"
	id := uuid.New().String()
	b.Event(wdgo.Cmd{"AddStage", id})
	s, _ := b.Find(id)
	s.Event(wdgo.Cmd{"Name", "Backlog"})
	s.Event(wdgo.Cmd{"AddCard", "c1"})
	c, _ := b.Find("c1")
	c.Event(wdgo.Cmd{"Name", "Todo 1"})
	s.Event(wdgo.Cmd{"AddCard", "c2"})
	c, _ = b.Find("c2")
	c.Event(wdgo.Cmd{"Name", "Todo 2"})
	s.Event(wdgo.Cmd{"AddCard", "c3"})
	c, _ = b.Find("c3")
	c.Event(wdgo.Cmd{"Name", "Todo 3"})

	id = uuid.New().String()
	b.Event(wdgo.Cmd{"AddStage", id})
	s, _ = b.Find(id)
	s.Event(wdgo.Cmd{"Name", "Doing"})

	id = uuid.New().String()
	b.Event(wdgo.Cmd{"AddStage", id})
	s, _ = b.Find(id)
	s.Event(wdgo.Cmd{"Name", "Done"})
	a.addBoard(b)
	b = wdgo.NewBoard(uuid.New().String())
	b.Name = "Board 2"
	a.addBoard(b)
}

func (a *app) addBoard(b *wdgo.Board) {
	a.boards = append(a.boards, b)
}

func (a *app) run() error {
	err := a.layout()
	if err != nil {
		return fmt.Errorf("app.run().layout(): %w", err)
	}
	return a.view.root.Run()
}

func (a *app) layout() error {
	a.initBoards()
	a.setHome()
	a.view.board = tview.NewFlex()
	a.view.card = tview.NewForm()
	//a.setBoard()
	//a.setCard()
	a.view.pages = tview.NewPages().
		AddPage("home", a.view.home, true, true)

	a.view.root = tview.NewApplication().SetRoot(a.view.pages, true).
		EnableMouse(true)
	a.view.root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyBackspace {
			a.view.pages.SwitchToPage("home")
			a.view.root.SetFocus(a.view.home)
		}
		return event
	})
	return nil
}

func (a *app) setHome() error {
	a.view.home = tview.NewList()
	for i, b := range a.boards {
		abc := "abcdefghijklmnopqrstuvwxyz"
		a.view.home.AddItem(b.Name, b.ID(), rune(abc[i%25]), func() {
			a.selectedBoard = a.boards[a.view.home.GetCurrentItem()]
			a.setBoard()
			//a.view.root.SetRoot(a.view.board, true)
			//a.view.pages.AddAndSwitchToPage("board", a.view.board, true)
			a.view.pages.SwitchToPage("board")
		})
	}
	a.view.home.AddItem("New Board", "create a new board", '+', nil)
	a.view.home.AddItem("Quit", "", 'Q', func() {
		a.view.root.Stop()
	})
	return nil
}
func (a *app) setBoard() error {
	if a.selectedBoard == nil {
		return fmt.Errorf("setBoard(): no board selected")
	}

	a.view.board = tview.NewFlex()
	for _, s := range a.selectedBoard.Stages {
		stage := tview.NewFlex()
		stage.SetTitle(s.Name).SetBorder(true)
		cards := tview.NewFlex().SetDirection(tview.FlexRow)
		cards.SetBorder(false)
		for _, c := range s.Cards {
			txt := tview.NewTextView()
			txt.SetBorder(true)
			io.WriteString(txt, c.Name)
			cards.AddItem(txt, 3, 1, false)
		}
		if len(s.Cards) > 0 {
			stage.AddItem(cards, 0, 1, false)
		}
		a.view.board.AddItem(stage, 0, 1, false)
	}
	a.view.pages.AddPage("board", a.view.board, true, false)

	return nil
}
func (a *app) setCard() error {
	a.view.card = tview.NewForm()
	a.view.card.SetTitle("Card")
	return nil
}
func (a *app) setFoo() error { return nil }
