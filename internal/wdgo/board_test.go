package wdgo

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestBoard_Event(t *testing.T) {
	initialBoard := NewBoard("1")
	initialBoard.Name = "empty"

	tests := []struct {
		name    string
		b       *Board
		cmd     Cmd
		want    string
		wantErr bool
	}{
		{
			"set name: foo",
			initialBoard,
			Cmd{"Name", "foo"},
			"Name:\"foo\"",
			false,
		},
		{
			"not existing action",
			initialBoard,
			Cmd{"AnotherName", "xxx"},
			"Name:\"foo\"",
			true,
		},
		{
			"set name: bar",
			initialBoard,
			Cmd{"Name", "bar"},
			"Name:\"bar\"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.Event(tt.cmd); (err != nil) != tt.wantErr {
				t.Errorf("Board.Event() error = %v, wantErr %v", err, tt.wantErr)
			}
			bStr := fmt.Sprintf("%#v", *tt.b)
			if !strings.Contains(bStr, tt.want) {
				t.Errorf("Got: %#v\nWant: %s", *tt.b, tt.want)
			}
		})
	}
}

func TestBoardFind(t *testing.T) {
	b := NewBoard("1")
	b.Name = "Testboard"
	b.Event(Cmd{"AddStage", "s0"})
	b.Event(Cmd{"AddStage", "s1"})
	b.Event(Cmd{"AddStage", "s2"})
	b.Event(Cmd{"AddStage", "s3"})
	f, err := b.Find("1")
	if !reflect.DeepEqual(f, b) {
		t.Errorf("expect to find the board")
	}
	expectNoError(t, err)
	f, err = b.Find("s1")
	if f != b.Stages[1] {
		t.Errorf("expect to find stage s1")
	}
	expectNoError(t, err)
	_, err = b.Find("nothing")
	if !errors.Is(err, ErrIDNotFound) {
		t.Errorf("expect IDNotFound. got: %s", err)
	}
	f, err = b.Find("s2")
	if f != b.Stages[2] {
		t.Errorf("expect to find stage s2")
	}
	expectNoError(t, err)
}

func expectNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("no error expected. got: %s", err)
	}
}
