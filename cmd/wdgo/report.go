package main

import (
	"fmt"
	"sort"
	"strconv"
	"time"

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

// renderReport creates the report
func (a *app) renderReport() {
	a.report.Clear()
	data := a.createSessionTable()

	for ri, row := range data {
		for ci, cell := range row {
			if ci == 0 {
				continue
			}
			align := tview.AlignRight
			switch ci {
			case 2, 3:
				align = tview.AlignCenter
			default:
				align = tview.AlignRight
			}
			a.report.SetCell(ri, ci,
				tview.NewTableCell(cell).SetAlign(align))
		}
	}
	a.report.SetSelectable(true, false)
	a.report.ScrollToBeginning()
	a.pages.AddAndSwitchToPage("report", a.report, true)
	a.root.SetFocus(a.report)
}

func (a *app) createSessionTable() [][]string {
	const (
		dayFormat  = "02.01.06"
		sortFormat = "06-01-02"
		timeFormat = "15:04"
	)
	data := [][]string{}
	for _, b := range a.boards {
		b.aggregator.State()
		for _, stage := range b.board.Stages {
			for _, card := range stage.Cards {
				var cardsum time.Duration
				for _, s := range card.Sessions {
					sDuration := s.End.Sub(s.Start)
					if sDuration > 0 {
						cardsum += sDuration
					}
					row := make([]string, 11)
					row[0] = fmt.Sprintf("%s:%s",
						s.Start.Format(sortFormat),
						card.ID())
					row[1] = s.Start.Format(dayFormat)
					row[2] = b.board.Name
					row[3] = card.Name
					row[4] = card.Description
					row[5] = card.SupportID
					row[6] = s.Start.Format(timeFormat)
					row[7] = s.End.Format(timeFormat)
					row[8] = fmt.Sprintf("%.2f", sDuration.Hours())
					if sDuration < 0 {
						row[8] = "-"
					}
					row[9] = fmt.Sprintf("%.2f", cardsum.Hours())
					data = append(data, row)
				}
			}
		}
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i][0] > data[j][0]
	})
	var dayCardSum float64
	var lastSortVal string
	for _, row := range data {
		if lastSortVal == "" || lastSortVal != row[0] {
			lastSortVal = row[0]
			dayCardSum = 0
		}
		dur, err := strconv.ParseFloat(row[8], 32)
		if err == nil {
			dayCardSum += dur
		}
		row[10] = fmt.Sprintf("%.2f", dayCardSum)
	}
	return data
}
