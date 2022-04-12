package core

func Controller(ansibleDir string, playbookDir string) {

	app := new(Config)
	app.AnsibleDir  = ansibleDir
	app.PlaybookDir = playbookDir
	app.Playbooks   = ParsePlaybooks(playbookDir)

	UI := &TUI{
	 Playbooks:   app.Playbooks,
	 PlaybookDir: app.PlaybookDir,
	 AnsibleDir:  app.AnsibleDir,
	}
	app.UI = UI.Run()

}

