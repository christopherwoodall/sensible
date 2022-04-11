package core

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)


func parse_header(playbook_path string) (*Playbook, bool) {
	playbook, err := ioutil.ReadFile(playbook_path); if err != nil {
		return nil, true
	}
	playbook_content := string(playbook)
	sensible_header  :="### sensible ###"
	sensible_footer  := "### /sensible ###"

	if !strings.Contains(playbook_content, sensible_header) ||
		 !strings.Contains(playbook_content, sensible_footer) {
		return nil, true
	}

	header := strings.Split(playbook_content, sensible_footer)[0]
	header  = strings.Split(header, sensible_header)[1]
	header  = strings.Replace(header, "# ", "", -1)

	playbook_header := new(Playbook)
	if err := yaml.Unmarshal([]byte(header), playbook_header); err != nil {
		// Invalid Sensible header format.
		switch {
		default:
		}
	}
	return playbook_header, false
}


func ParsePlaybooks(playbook_dir string) []Playbook {
	glob := filepath.Join(playbook_dir, "/*.*")
	files, _ := filepath.Glob(glob)
	playbooks := []Playbook{}

	for _, file := range files {
		switch filepath.Ext(file) {
			case
				".yaml",
				".yml":
					header, err := parse_header(file)
					if !err {
						playbooks = append(playbooks, *header)
					}
			default:
				continue
		}
	}
	return playbooks
}
