package core

import (
	"encoding/json"
	"fmt"
)


func Controller(ansibleDir string, playbookDir string) {

	app := new(Config)
	app.AnsibleDir  = ansibleDir
	app.PlaybookDir = playbookDir
	app.Playbooks   = ParsePlaybooks(playbookDir)

	app_debug, _ := json.MarshalIndent(app, "", "\t")
	fmt.Print(string(app_debug))

	// app.UI = New()
	//Start_TUI(app.Playbooks)
	UI := &TUI{
	 Playbooks: app.Playbooks,
	}
	app.UI = UI.Run()

}
