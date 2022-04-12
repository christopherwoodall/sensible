package core

import "github.com/rivo/tview"

type Config struct {
	AnsibleDir string
	PlaybookDir string
	Playbooks []Playbook
	UI *TUI
}


type Playbook struct {
	Name        string   			 `yaml:"name"`
	Description string   			 `yaml:"description"`
	Index       int      			 `yaml:"index"`
	Tags 				[]string	     `yaml:"tags"`
	Args        []string 	     `yaml:"args"`
	Requires    []string			 `yaml:"requires"`
	Vars				[]PlaybookVars `yaml:"vars"`
	Path        string
	Selected    bool
}

type PlaybookVars struct {
	Key string
	Value string
}


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

