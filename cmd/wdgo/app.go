package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/as27/wdgo/internal/estore"
	"github.com/as27/wdgo/internal/wdgo"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type app struct {
	root        *tview.Application
	pages       *tview.Pages
	home        *tview.List
	card        *tview.Form
	activeBoard int
	boards      []board
	path        appPaths
}

type board struct {
	view        *tview.Flex
	board       *wdgo.Board
	aggregator  *estore.Aggregator
	stageviews  []*tview.Flex
	cardviews   [][]*tview.TextView
	activeStage int
	activeCard  int
}

type appPaths struct {
	app   string
	event string
}

func newApp(p appPaths) *app {
	a := app{
		root:  tview.NewApplication(),
		pages: tview.NewPages(),
		home:  tview.NewList(),
		card:  tview.NewForm(),
		path:  p,
	}
	err := a.initBoards()
	if err != nil {
		log.Println("cannot init boards")
		panic(err)
	}
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
		a.stop()
	}
	return event
}

func (a *app) setCard() error {
	//	a.view.card = tview.NewForm()
	//	a.view.card.SetTitle("Card")
	return nil
}

func (a *app) addBoard(b *wdgo.Board) {
	bb := board{
		view:        tview.NewFlex(),
		board:       b,
		aggregator:  estore.NewAggregator(b),
		stageviews:  []*tview.Flex{},
		cardviews:   [][]*tview.TextView{},
		activeStage: 0,
		activeCard:  0,
	}
	a.boards = append(a.boards, bb)
}

func (a *app) writeBoards(w io.Writer) error {
	for _, b := range a.boards {
		fmt.Fprintf(w, "%s|%s\n", b.board.ID(), b.board.Name)
	}
	return nil
}

// readBoards expect a simple logic
// [id]|[name]
// 12312323-21232a|board 1
func (a *app) readBoards(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		txt := scanner.Text()
		e := strings.Split(txt, "|")
		if len(e) != 2 {
			return fmt.Errorf("wrong input line: %s", txt)
		}
		id, name := e[0], e[1]
		b := wdgo.NewBoard(id)
		b.Name = name
		a.addBoard(b)
	}
	return nil
}

func (a *app) initBoards() error {
	bfile, err := os.Open(a.path.app)
	if err != nil {
		return fmt.Errorf("app.initBoards(): %w", err)
	}
	defer bfile.Close()
	err = a.readBoards(bfile)
	if err != nil {
		return fmt.Errorf("app.stop().initBoards.read: %w", err)
	}
	return nil
}

func (a *app) stop() error {
	bfile, err := os.Create(a.path.app)
	if err != nil {
		return fmt.Errorf("app.stop(): %w", err)
	}
	defer bfile.Close()
	err = a.writeBoards(bfile)
	if err != nil {
		return fmt.Errorf("app.stop().writeBoards: %w", err)
	}
	a.root.Stop()
	return nil
}
