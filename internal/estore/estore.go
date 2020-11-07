// Package estore handles events from different sources
// It safes the events on a source and it sends the events
// to the model to aggregate the state
package estore

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/as27/wdgo/internal/wdgo"
)

const fileSeparator = "|"

// Event is used by the aggregator to create a state of the
// board.
type Event struct {
	id   string
	time time.Time
	cmd  wdgo.Cmd
}

func NewEvent(id, action, value string) Event {
	return Event{
		id:   id,
		time: time.Now(),
		cmd: wdgo.Cmd{
			Action: action,
			Value:  value,
		},
	}
}

// Aggregator is used to aggregate the state of a board over
// events
type Aggregator struct {
	board   *wdgo.Board
	events  []Event
	version int
}

// NewAggregator creates an empty new Aggregator. The state
// will be written inside the board.
func NewAggregator(board *wdgo.Board) *Aggregator {
	return &Aggregator{
		board:   board,
		events:  []Event{},
		version: 0,
	}
}

// Event adds an event to the end of the aggregate
func (a *Aggregator) Event(e Event) {
	a.events = append(a.events, e)
}

func (a *Aggregator) NewEvent(id, action, value string) {
	a.Event(NewEvent(id, action, value))
}

// Init deletes the state of the board.
func (a *Aggregator) Init() {
	//a.board.Name = ""
	a.board.Created = time.Time{}
	a.board.Modified = time.Time{}
	a.board.Stages = []*wdgo.Stage{}
}

// State aggregates all events until the end
func (a *Aggregator) State() {
	a.Init()
	a.Version(len(a.events))
}

// Version aggregates the state until a specific version
func (a *Aggregator) Version(version int) {
	a.Init()
	for i, e := range a.events {
		if i > version {
			break
		}
		ee, err := a.board.Find(e.id)
		if err != nil {
			log.Println("Error: Aggregator.State: cannot find id: ", err)
			continue
		}
		err = ee.Event(e.cmd)
		if err != nil {
			log.Println("Error: Aggregator.State: ee.Event: ", err)
		}
		ee.Action(e.time)
		a.version = i
	}
}

func (a *Aggregator) SaveEvents(w io.Writer) error {
	for _, e := range a.events {
		s := []string{
			e.time.Format(wdgo.TimeFormat),
			e.id,
			e.cmd.Action,
			e.cmd.Value,
		}
		_, err := fmt.Fprintln(w, strings.Join(s, fileSeparator))
		if err != nil {
			return fmt.Errorf("aggregator.SaveEvents: %w", err)
		}
	}
	return nil
}

func (a *Aggregator) LoadEvents(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), fileSeparator)
		t, err := time.Parse(wdgo.TimeFormat, s[0])
		if err != nil {
			return fmt.Errorf("aggregator.LoadEvents: %w", err)
		}
		e := Event{
			time: t,
			id:   s[1],
			cmd: wdgo.Cmd{
				Action: s[2],
				Value:  s[3],
			},
		}
		a.Event(e)
	}
	return nil
}
