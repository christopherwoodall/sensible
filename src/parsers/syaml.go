package parsers

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Header struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Index       int      `yaml:"index"`
	Args        []string `yaml:"args"`
	Requires    []string `yaml:"requires"`
	Tags        []string `yaml:"tags"`
	Selected    bool     `yaml:"selected"`
	Path        string   `yaml:"path"`
}

func check(e error) {
	if e != nil {
		fmt.Println("Error:", e.Error())
		os.Exit(1)
	}
}

func clean_array(headers []Header) []Header {
	var cleaned []Header
	for _, h := range headers {
		if h.Name != "" {
			cleaned = append(cleaned, h)
		}
	}
	return cleaned
}

func syaml(playbook_path string) (Header, error) {
	header := new(Header)
	playbook, err := ioutil.ReadFile(playbook_path)
	check(err)

	playbook_str := string(playbook)
	playbook_str = strings.Split(playbook_str, "### /sensible ###")[0]
	playbook_str = strings.Split(playbook_str, "### sensible ###")[1]
	playbook_str = strings.Replace(playbook_str, "# ", "", -1)

	if err = yaml.Unmarshal([]byte(playbook_str), header); err != nil {
		// Invalid Sensible header format.
		return *header, err
	}
	header.Path = playbook_path

	return *header, nil
}

func Playbooks(playbooks []string) []Header {
	// headers := []Header{}
	headers := make([]Header, 1000)

	for _, file := range playbooks {
		header, err := syaml(file)
		if err != nil {
			continue
		}
		// headers = append(headers, header)
		headers[header.Index] = header
	}

	return clean_array(headers)
}
