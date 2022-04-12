package core

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TUI struct {
	App *tview.Application
	Header *tview.TextView
	Menu *tview.Table
	Details *tview.TextView
	Chyron *tview.Table
	// Modal
	PlaybookDir string
	AnsibleDir string
	Playbooks []Playbook
	CurIndex int
}


func (tui *TUI) Run() *TUI {
	tui.App    = tview.NewApplication()

	tui.Header = func() *tview.TextView {
		return tview.NewTextView().
								 SetTextAlign(tview.AlignCenter).
								 SetText("Sensible")
	}()

	tui.Menu = tview.NewTable()
	tui.Menu.
		SetBorders(false).
		SetSelectable(true, true)

	tui.Details = tview.NewTextView()
	tui.Details.
		SetWrap(false).
		SetDynamicColors(true).
		SetBorderPadding(1, 1, 2, 0)

		tui.Chyron = tui.DrawFooter()

	tui.Draw()

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

	tui.globalEventHanbler()

	if err := tui.App.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	return tui
}


/////////////////////////////////////////
//
func (tui *TUI) globalEventHanbler() {
	table := tui.Menu
	tui.Menu.Select(0, 0).
	SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape { tui.App.Stop()     }
		if key == tcell.KeyEnter  { tui.run_selected() }
	}).
	SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).
			SetTextColor(tcell.ColorBlue)
	}).
	SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		tui.Draw()

		switch event.Key() {
			case tcell.KeyEnter:
				tui.run_selected()
			case tcell.KeyUp:
				if (tui.CurIndex - 1) > 0 {
					tui.CurIndex = tui.CurIndex - 1
				} else {
					tui.CurIndex = 0
				}
			case tcell.KeyDown:
				if (tui.CurIndex + 1) < len(tui.Playbooks) {
					tui.CurIndex = tui.CurIndex + 1
				} else {
					tui.CurIndex = len(tui.Playbooks) - 1
				}
			// TODO: Mouse Handling
			default:
				_, tui.CurIndex = table.GetSelection()
		}

		switch event.Rune() {
			case ' ':
				tui.mark_selected()
		}
		//tui.highlight_node()
		return event
	})
	tui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
			case 'q':
				tui.App.Stop()
		}
		return event
	})
}


/////////////////////////////////////////
// TODO: Better Path management
func (tui *TUI) run_selected() {
	tui.App.Stop()

	os.Chdir(tui.AnsibleDir)
	for _, playbook := range tui.Playbooks {
		ansible_cmd := "ansible-playbook"
		if playbook.Selected {
			ansible_cmd += " playbooks/" + playbook.Path
			for _, arg := range playbook.Args {
				ansible_cmd += " " + arg
			}
			ansible_cmd += "</dev/tty >/dev/tty 2>/dev/tty"
			fmt.Println(ansible_cmd)
			// os.StartProcess(ansible_cmd, []string{}, &os.ProcAttr{Dir: tui.AnsibleDir})
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			cmd := exec.Command("bash", "-c", ansible_cmd)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			cmd.Run()
			fmt.Println(stdout.String())
		}
	}
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
func (tui *TUI) Draw() {
	tui.Menu.Clear()
	tui.Details.Clear()

	tui.Menu    = tui.DrawMenu()
	tui.Details = tui.DrawDetails()
}


/////////////////////////////////////////
//
func (tui *TUI) mark_selected() {
	// row, col := tui.CurIndex, 0
	row, col := tui.Menu.GetSelection()
	playbook := tui.Playbooks[row]
	tui.Playbooks[row].Selected = ! tui.Playbooks[row].Selected
	if tag_in(playbook.Tags, "seperator") { return }
	if tui.Playbooks[row].Selected {
		tui.Menu.GetCell(row, col).
			SetTextColor(tcell.ColorBlue).
			SetText("  [✓] " + playbook.Name)
	} else {
		tui.Menu.GetCell(row, col).
			SetTextColor(tcell.ColorWhite).
			SetText("  [ ] " + playbook.Name)
	}
}

func (tui *TUI) highlight_node() {
	// row, col := tui.CurIndex, 0
	row, col := tui.Menu.GetSelection()
	playbook := tui.Playbooks[row]
	if ! tag_in(playbook.Tags, "seperator") {
		content := tui.Menu.GetCell(row, col).Text
		content = ">" + content[1:]
		tui.Menu.GetCell(row, col).
			SetText(content)
	}
}


/////////////////////////////////////////
//
func (tui *TUI) DrawMenu() *tview.Table {
	table := tui.Menu
	_, _, width, _ := table.GetInnerRect()

	for i, playbook := range tui.Playbooks {
		var cell *tview.TableCell
		var content string
		color := tcell.ColorWhite
		switch  {
			case
				tag_in(playbook.Tags, "seperator"):
					padding := strings.Repeat("=", width)
					color    = tcell.Color51
					content  = padding + " " + playbook.Name + " " + padding
					cell = tview.
						NewTableCell(content).
						SetTextColor(color).
						SetAlign(tview.AlignCenter).
						SetExpansion(1)
			default:
				if playbook.Selected {
					color   = tcell.ColorBlue
					content = "  [✓] " + playbook.Name
				} else{
					content = "  [ ] " +  playbook.Name
				}
				cell = tview.
					NewTableCell(content).
					SetTextColor(color).
					SetAlign(tview.AlignLeft).
					SetExpansion(1)
		}
		table.SetCell(i, 0, cell)
	}
	return table
}


/////////////////////////////////////////
//
func (tui *TUI) DrawDetails() *tview.TextView {
	row, _ := tui.Menu.GetSelection()
	current_playbook := tui.Playbooks[row]
	tui.Details.SetText(current_playbook.Description)
	return tui.Details
}


/////////////////////////////////////////
//
func (tui *TUI) DrawFooter() *tview.Table {
	footer := tview.NewTable()
	footer.
		SetBorders(false).
		SetSelectable(false, false)
	help_txt := []string{
		"q: Quit",
		"Space: Select",
		"Enter: Run",
	}
	for i, line := range help_txt {
		footer.SetCell(0, i, tview.NewTableCell(line + " | ").
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignCenter))
	}
	return footer
}


/////////////////////////////////////////
// NOTES
//  - https://github.com/rivo/tview/blob/master/table.go
//  -

