///usr/bin/env -S go run ./ "$@"; exit "$?"

package main

import (
	"flag"
	"fmt"
	"os"
	core "sensible/core"
)

const (
	AUTHOR  = "Christopher Woodall"
	NAME    = "Sensible"
	VERSION = "0.0.1"
)


func main() {
	ansibleDir := func() string {
		var target_dir string
		flag.StringVar(
			&target_dir,
			"dir",
			"../playbooks",
			"The target directory containing the playbooks")
		flag.Parse()
		if _, err := os.Stat(target_dir); os.IsNotExist(err) {
			fmt.Println("[!] Target directory does not exist")
			os.Exit(1)
		}
		return target_dir
	}()

	playbookDir, err := FindPlaybookDirectory(ansibleDir); if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	core.Controller(ansibleDir, playbookDir)
}

