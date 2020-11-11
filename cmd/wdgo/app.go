package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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
	card        cardViews
	stage       *tview.Form
	newBoard    *tview.Form
	activeBoard int
	boards      []board
	path        appPaths
}

type board struct {
	view         *tview.Flex
	board        *wdgo.Board
	aggregator   *estore.Aggregator
	stageviews   []*tview.Flex
	cardviews    [][]*tview.TextView
	activeStage  int
	activeCard   int
	cardSelected *wdgo.Card
}

type cardViews struct {
	card         *tview.Flex
	form         *tview.Form
	sessionsFlex *tview.Flex
	sessions     *tview.Table
	sessionForm  *tview.Form
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
		card: cardViews{
			card:         tview.NewFlex(),
			form:         tview.NewForm(),
			sessionsFlex: tview.NewFlex(),
			sessions:     tview.NewTable(),
			sessionForm:  tview.NewForm(),
		},
		stage:    tview.NewForm(),
		newBoard: tview.NewForm(),
		path:     p,
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
		a.sessionEvents(event)
	case "stage":
		a.stageEvents(event)
	}
	// global key bindings
	if event.Key() == tcell.KeyCtrlQ {
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
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("app.initBoards(): %w", err)
	}
	defer bfile.Close()
	err = a.readBoards(bfile)
	if err != nil {
		return fmt.Errorf("app.stop().initBoards.read: %w", err)
	}
	err = a.loadEvents()
	if err != nil {
		return fmt.Errorf("app.stop().loadEvents.read: %w", err)
	}
	return nil
}

func (a *app) stop() error {
	err := os.MkdirAll(filepath.Dir(a.path.app), 0666)
	if err != nil {
		return fmt.Errorf("app.stop():MkdirAll: %w", err)
	}
	bfile, err := os.Create(a.path.app)
	if err != nil {
		return fmt.Errorf("app.stop(): %w", err)
	}
	defer bfile.Close()
	err = a.writeBoards(bfile)
	if err != nil {
		return fmt.Errorf("app.stop().writeBoards: %w", err)
	}
	err = a.writeEvents()
	if err != nil {
		return fmt.Errorf("app.stop().writeEvents: %w", err)
	}
	a.root.Stop()
	return nil
}

func (a *app) writeEvents() error {
	err := os.MkdirAll(a.path.event, 0666)
	if err != nil {
		return fmt.Errorf("app.writeEvents() MkdirAll: %w", err)
	}
	for _, b := range a.boards {
		fd, err := os.Create(filepath.Join(a.path.event, b.board.ID()))
		if err != nil {
			return fmt.Errorf("writeEvents.Create:%s.:%w", b.board.Name, err)
		}
		err = b.aggregator.SaveEvents(fd)
		if err != nil {
			return fmt.Errorf("writeEvents.%s.:%w", b.board.Name, err)
		}
		fd.Close()
	}
	return nil
}

func (a *app) loadEvents() error {
	for _, b := range a.boards {
		fd, err := os.Open(filepath.Join(a.path.event, b.board.ID()))
		if err != nil {
			return fmt.Errorf("loadEvents.%s.:%w", b.board.Name, err)
		}
		err = b.aggregator.LoadEvents(fd)
		if err != nil {
			return fmt.Errorf("loadEvents.%s.:%w", b.board.Name, err)
		}
		fd.Close()
		b.aggregator.State()
	}
	return nil
}
