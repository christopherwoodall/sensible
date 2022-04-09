package main

import (
	"fmt"
	"os"

	"sensible/parsers"
	"sensible/utils"
)

func check(e error) {
	if e != nil {
		fmt.Println("Error:", e.Error())
		os.Exit(1)
	}
}

func main() {
	target_dir, err := utils.ArgParse()
	check(err)

	playbooks, err := utils.ScanDir(target_dir)
	check(err)

	headers := parsers.Playbooks(playbooks)
	if len(headers) == 0 {
		fmt.Println("No playbooks found.")
		os.Exit(1)
	}

	// fmt.Println(target_dir)
	// fmt.Println(playbooks)
	// for _, h := range headers {
	// 	fmt.Println(h)
	// }

	TUI(headers)
}
