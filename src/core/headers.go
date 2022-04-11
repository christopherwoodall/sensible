package core

type Config struct {
	AnsibleDir string
	PlaybookDir string
	Playbooks []Playbook
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
