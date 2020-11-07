package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestAppLoadBoard(t *testing.T) {
}

func TestReadBoardsWriteBoards(t *testing.T) {
	tests := []struct {
		name      string
		a         *app
		r         io.Reader
		wantErr   bool
		wantIDs   []string
		wantNames []string
	}{
		{
			"one entry",
			&app{},
			strings.NewReader("id1|board 1\n"),
			false,
			[]string{"id1"},
			[]string{"board 1"},
		},
		{
			"more entries",
			&app{},
			strings.NewReader("id1|b 1\nid2|b 2\nid3|b 3\n"),
			false,
			[]string{"id1", "id2", "id3"},
			[]string{"b 1", "b 2", "b 3"},
		},
		{
			"wrong line",
			&app{},
			strings.NewReader("id1 2\nid3|b 3\n"),
			true,
			[]string{"id1", "id2", "id3"},
			[]string{"b 1", "b 2", "b 3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inBuffer := &bytes.Buffer{}
			r := io.TeeReader(tt.r, inBuffer)
			err := tt.a.readBoards(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("app.loadBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i, b := range tt.a.boards {
				if tt.wantIDs[i] != b.board.ID() {
					t.Errorf("app.loadBoard() wrong id: %s, want: %s",
						b.board.ID(), tt.wantIDs[i])
				}
				if tt.wantNames[i] != b.board.Name {
					t.Errorf("app.loadBoard() wrong name: %s, want: %s",
						b.board.Name, tt.wantNames[i])
				}

				t.Run("writeBoard", func(t *testing.T) {
					if tt.wantErr {
						// don't write, when read-test expects
						// an error
						return
					}
					buf := &bytes.Buffer{}
					tt.a.writeBoards(buf)
					got := buf.String()
					wantIn := inBuffer.String()
					if got != wantIn {
						t.Errorf("writeBoard()=\n%s\nwant:\n%s", got, wantIn)
					}
				})
			}
		})
	}
}
