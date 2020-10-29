package wdgo

import (
	"errors"
	"time"
)

var (
	ErrNoSuchAction = errors.New("no such a action")
	ErrIDNotFound   = errors.New("no such ID")
)

const (
	TimeFormat = "2006-01-02 15:04:05"
)

type Eventer interface {
	Event(Cmd) error
}

type Cmd struct {
	Action string
	Value  string
}

type EventSource struct {
	ID       string
	Name     string
	Created  time.Time
	Modified time.Time
}

func (e EventSource) GetID() string {
	return e.ID
}
