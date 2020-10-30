package wdgo

import (
	"errors"
	"time"
)

// Predefined errors for the wdgo-package.
// ErrNoSuchAction is used for the Eventer-Interface
// ErrIDNotFound is used for the Find-methods
var (
	ErrNoSuchAction = errors.New("no such a action")
	ErrIDNotFound   = errors.New("no such ID")
)

// TimeFormat defines the timestamp for the package
const (
	TimeFormat = "2006-01-02 15:04:05"
)

// Eventer is the interface to abstract the types for this
// package. The cmd includes the action and a value, which is
// used for the event.
// If a type does not support the given action ErrIDNotFound has
// to be returned.
type Eventer interface {
	Event(Cmd) error
}

// Cmd defines the structure a event is defined. The Eventer needs
// to implement every action. The values is set by a string. So all
// values need to be stringyfied.
type Cmd struct {
	Action string
	Value  string
}

// EventSource is the basic type for an eventer.
type EventSource struct {
	ID       string
	Name     string
	Created  time.Time
	Modified time.Time
}

// Action is used to set the timestamps inside the EventSource
// Calling that method the first time, also Created timestamp
// is set. The caller defines the time for the timestamp.
func (e *EventSource) Action(t time.Time) {
	if (e.Created == time.Time{}) {
		e.Created = t
	}
	e.Modified = t
}
