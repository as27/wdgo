package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

func (a *app) homeEvents(event *tcell.EventKey) *tcell.EventKey {
	return event
}

func (a *app) renderHome() error {
	a.home.Clear()
	for i, b := range a.boards {
		abc := "abcdefghijklmnopqrstuvwxyz"
		a.home.AddItem(b.board.Name, b.board.ID(), rune(abc[i%25]),
			func() {
				a.activeBoard = a.home.GetCurrentItem()
				a.renderBoard()
				//a.pages.SwitchToPage("board")
				//a.root.SetFocus(a.pages)
			})
	}
	a.home.AddItem("New Board", "create a new board", '+', func() {
		a.renderNewBoard()
	})
	a.home.AddItem("Quit", "", 'Q', func() {
		err := a.stop()
		if err != nil {
			log.Println(err)
		}
	})
	a.pages.AddAndSwitchToPage("home", a.home, true)

	return nil
}
