package utils

import (
	"errors"
	"flag"
	"os"
)

func ArgParse() (string, error) {

	var target_dir string
	// var files []string

	flag.StringVar(
		&target_dir,
		"dir",
		"../playbooks",
		"The target directory containing the playbooks")

	flag.Parse()

	if _, err := os.Stat(target_dir); os.IsNotExist(err) {
		return "", errors.New("target directory does not exist")
	}

	return target_dir, nil

}
