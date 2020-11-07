// Package estore handles events from different sources
// It safes the events on a source and it sends the events
// to the model to aggregate the state
package estore

import (
	"log"
	"time"

	"github.com/as27/wdgo/internal/wdgo"
)

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
