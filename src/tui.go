///usr/bin/env -S go run "$0" "$@"; exit "$?"
package main

import (
	"sensible/parsers"

	"github.com/rivo/tview"
)

func create_list(items []parsers.Header, app *tview.Application) *tview.List {
	list := tview.NewList()

	for i, item := range items {
		// rune := rune(i + 65)
		list.AddItem(item.Name, "", rune(i), nil)
	}
	list.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})
	return list
}

func TUI(playbooks []parsers.Header) {
	app := tview.NewApplication()

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	header := newPrimitive("Sensible")
	chyron := newPrimitive("Footer")

	menu := newPrimitive("Menu")
	main := newPrimitive("Main content")
	sideBar := newPrimitive("Side Bar")

	// list := tview.NewList().
	// 	AddItem("List item 1", "Some explanatory text", 'a', nil).
	// 	AddItem("List item 2", "Some explanatory text", 'b', nil).
	// 	AddItem("List item 3", "Some explanatory text", 'c', nil).
	// 	AddItem("List item 4", "Some explanatory text", 'd', nil).
	// 	AddItem("Quit", "Press to exit", 'q', func() {
	// 		app.Stop()
	// 	})
	list := create_list(playbooks, app)

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(header, 0, 0, 1, 3, 0, 0, false).
		AddItem(chyron, 2, 0, 1, 3, 0, 0, false)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(menu, 0, 0, 0, 0, 0, 0, false).
		// AddItem(main, 1, 0, 1, 3, 0, 0, false).
		AddItem(list, 1, 0, 1, 3, 0, 0, false).
		AddItem(sideBar, 0, 0, 0, 0, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, false).
		AddItem(list, 1, 1, 1, 1, 0, 100, false).
		AddItem(main, 1, 2, 1, 1, 0, 100, false)

	//if err := tview.NewApplication().SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
