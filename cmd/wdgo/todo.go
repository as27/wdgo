package main

import (
	"sort"

	"github.com/rivo/tview"
)

func (a *app) renderToDo() {
	a.report.Clear()

	data := a.createToDoTable()

	for ri, row := range data {
		for ci, cell := range row {
			if ci <= 1 {
				continue
			}
			align := tview.AlignLeft

			a.report.SetCell(ri, ci,
				tview.NewTableCell(cell).SetAlign(align))
		}
	}
	a.report.SetSelectable(true, false)
	a.report.SetSelectedFunc(func(row, col int) {
		// jump into the card view when a row is selected
		cardID := data[row][0]
		a.setActiveCard(cardID)
		a.renderCard("edit")
	})
	a.report.ScrollToBeginning()
	a.pages.AddAndSwitchToPage("report", a.report, true)
	a.root.SetFocus(a.report)
}

func (a *app) createToDoTable() [][]string {
	data := [][]string{}
	for _, b := range a.boards {
		b.aggregator.State()
		for _, stage := range b.board.Stages {
			for _, card := range stage.Cards {
				if !card.ToDo {
					continue
				}
				row := make([]string, 6)
				row[0] = card.ID()
				row[1] = stage.Name
				row[2] = b.board.Name
				row[3] = card.Name
				row[4] = card.Description
				row[5] = card.SupportID
				data = append(data, row)

			}
		}
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i][0] > data[j][0]
	})
	return data
}
