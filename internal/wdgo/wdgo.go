package wdgo

import (
	"errors"
	"fmt"
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

// EventerID takes an eventer and returns the ID. The function
// uses the GetID Method of the EventSource.
func EventerID(e Eventer) (string, error) {
	type ider interface {
		GetID() string
	}
	i, ok := e.(ider)
	if !ok {
		return "", fmt.Errorf("%T not supported", e)
	}
	return i.GetID(), nil
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

func (e EventSource) GetID() string {
	return e.ID
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
