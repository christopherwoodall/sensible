package core

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)


type TUI struct {
	App *tview.Application
	Header *tview.TextView
	Menu *tview.List
	Details *tview.TextView
	Chyron *tview.TextView
	// Modal
	Playbooks []Playbook
	PrevIndex int
}


// func (tui *TUI) New() *TUI {
//   return &TUI{
// 		App: tview.NewApplication(),
// 		Header: tview.NewTextView(),
// 		Menu: tview.NewList(),
// 		Details: tview.NewTextView(),
// 		Chyron: tview.NewTextView(),
// 		Playbooks: []Playbook{},
// 		PrevIndex: 0,
// 	}
// }


func (tui *TUI) Run() *TUI {
	tui.App    = tview.NewApplication()
	tui.Header = func() *tview.TextView {
		return tview.NewTextView().
								 SetTextAlign(tview.AlignCenter).
								 SetText("Sensible")
	}()
	tui.Chyron = func() *tview.TextView {
		return tview.NewTextView().
								 SetTextAlign(tview.AlignCenter).
								 SetText("FOOTER")
	}()

	tui.Menu   = tview.NewList()
	tui.Details = tview.NewTextView()

	tui.Redraw()

	grid := tview.NewGrid().
								SetRows(1, 0, 3).
								SetColumns(-3, -3, -3).
								SetBorders(true)

	grid.AddItem(tui.Header, 0, 0, 1, 3, 0, 0, false).
			AddItem(tui.Chyron, 2, 0, 1, 3, 0, 0, false)
	// Layout for screens narrower than 100 cells (side bar is hidden)
	grid.AddItem(tui.Menu, 1, 0, 1, 3, 0, 0, true).
			 AddItem(tui.Details, 0, 0, 0, 0, 0, 0, false)
	// Layout for screens wider than 100 cells.
	grid.AddItem(tui.Menu, 1, 0, 1, 2, 0, 100, true).
			 AddItem(tui.Details, 1, 2, 1, 1, 0, 100, false)

	// tui.Listen()

	if err := tui.App.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	return tui
}


/////////////////////////////////////////
//
func tag_in(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}


/////////////////////////////////////////
//
func (tui *TUI) highlight_node() {
	curr_index := tui.Menu.GetCurrentItem()
	prev_index := tui.PrevIndex

	prev_selection, _ := tui.Menu.GetItemText(prev_index)
	curr_selection, _ := tui.Menu.GetItemText(curr_index)

	prev_content := strings.Replace(prev_selection, ">", " ", 1)
	curr_content := ">  " + curr_selection


	tui.Menu.SetItemText(prev_index, prev_content, "")
	tui.Menu.SetItemText(curr_index, curr_content, "")

}


/////////////////////////////////////////
// Event Listeners
func (tui *TUI) Listen( ) {
	tui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		curr_index := tui.Menu.GetCurrentItem()

		// tui.highlight_node()

		switch event.Rune() {
			case 'x':
				tui.App.Stop()
			case 'j':
				tui.Menu.SetItemText(curr_index, ">  ", "")
		}
		// if event.Key() == tcell.KeyCtrlL {
		// 	fmt.Println(tui.PrevIndex)
		// 	return nil
		// }
		tui.PrevIndex = curr_index
		return event
	})
}


/////////////////////////////////////////
//
func (tui *TUI) Redraw( ) {
	menu    := tui.draw_menu()
	details := func() *tview.TextView {
		return tview.NewTextView().
							   SetTextAlign(tview.AlignCenter).
	 							 SetText("DETAILS")
	}()

	tui.Menu    = menu
	tui.Details = details
}

func (tui *TUI) Quit( ) {
	tui.App.Stop()
}


/////////////////////////////////////////
//
func (tui *TUI) draw_menu( ) *tview.List {
	list := tui.Menu
	list.Clear()
	list.
		ShowSecondaryText(false).
		SetHighlightFullLine(true)

	for i, item := range tui.Playbooks {
		content := ""
		switch  {
			case
				tag_in(item.Tags, "seperator"):
				  _, _, width, _ := list.GetInnerRect()
					// margin  := strings.Repeat(" ", width / 2 )
					padding := strings.Repeat("=", width)
					// content = margin + padding + " " + item.Name + " " + padding
					content = padding + " " + item.Name + " " + padding
					list.AddItem(content, "", 0, nil)
			default:
				if tui.Playbooks[i].Selected {
					content = "  [X] " + item.Name
				} else {
					content = "  [ ] " + item.Name
				}
				list.AddItem(content, "", 0, func() {
					curr_index := list.GetCurrentItem()
					// curr_selection, _ := list.GetItemText(curr_index)
					// curr_content := strings.Replace(curr_selection, "[ ]", "[X]", 1)
					tui.Playbooks[curr_index].Selected = true
					// list.SetItemText(curr_index, curr_content, "")
					tui.Redraw()
				})
		}

	}

	list.AddItem("", "", 0, nil)
	list.AddItem("  Quit", "Press to exit", 0, tui.Quit)

	return list
}

