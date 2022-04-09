package utils

import (
	"errors"
	"os"
	"path/filepath"
)

func path_exist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func scan_directory(directory string) ([]string, error) {
	directory = filepath.Join(directory)
	types := []string{"yaml", "yml"}
	var discovered []string
	for _, t := range types {
		files, err := filepath.Glob(directory + "/*." + t)
		if err != nil {
			return nil, err
		}
		if len(files) > 0 {
			discovered = append(discovered, files...)
		}
	}
	return discovered, nil
}

func ScanDir(target_dir string) ([]string, error) {
	var files []string

	_paths := []string{"", "playbooks"}
	for _, p := range _paths {
		path := filepath.Join(target_dir, p)
		if path_exist(path) {
			found, err := scan_directory(path)
			if err != nil {
				return nil, errors.New("error scanning directory")
			}
			if len(found) > 0 {
				files = append(files, found...)
			}
		}
	}
	if len(files) == 0 {
		return nil, errors.New("no playbooks found")
	}

	return files, nil
}
