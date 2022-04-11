package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)


func path_exist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}


func check_for_header(playbook_path string) bool {
	playbook, _      := ioutil.ReadFile(playbook_path)
	playbook_content := string(playbook)
	sensible_header  :="### sensible ###"
	sensible_footer  := "### /sensible ###"

	if strings.Contains(playbook_content, sensible_header) ||
		 strings.Contains(playbook_content, sensible_footer) {
		return true
	}
	return false
}


func FindPlaybookDirectory(target_dir string) (string, error) {
	paths := []string{"", "playbooks"}

	if !path_exist(target_dir) {
		return "", errors.New("error scanning directory")
	}

	for _, path := range paths {
		path := filepath.Join(target_dir, path)
		glob := filepath.Join(path, "/*.*")
		if path_exist(path) {
			files, _ := filepath.Glob(glob)
			for _, file := range files {
				// if filepath.Base(file) == "ansible.cfg"
				switch filepath.Ext(file) {
				case
					".yaml",
					".yml":
						if check_for_header(file) {
							return path, nil
						}
				default:
					continue
				}
			}
		}
	}
	return "", errors.New("no playbook directory found")
}
