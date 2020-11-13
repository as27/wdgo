package main

import (
	"fmt"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *app) reportEvents(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyESC:
		a.pages.SwitchToPage("home")
		a.root.SetFocus(a.home)
	}
	return event
}

func (a *app) renderReport() {
	const (
		startFormat = "06-01-02 15:04"
		endFormat   = "15:04"
	)
	data := [][]string{}
	a.report.Clear()
	for _, b := range a.boards {
		b.aggregator.State()
		for _, stage := range b.board.Stages {
			for _, card := range stage.Cards {
				for _, s := range card.Sessions {
					row := make([]string, 10)
					row[0] = b.board.Name
					row[1] = card.Name
					row[2] = card.Description
					row[3] = card.SupportID
					row[4] = s.Start.Format(startFormat)
					row[5] = s.End.Format(endFormat)
					row[6] = fmt.Sprintf("%.2f", s.End.Sub(s.Start).Hours())
					data = append(data, row)
				}
			}
		}
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i][4] > data[j][4]
	})
	for ri, row := range data {
		for ci, cell := range row {
			a.report.SetCell(ri, ci, tview.NewTableCell(cell))
		}
	}
	a.report.SetSelectable(true, false)
	a.report.ScrollToBeginning()
	a.pages.AddAndSwitchToPage("report", a.report, true)
	a.root.SetFocus(a.report)
}
