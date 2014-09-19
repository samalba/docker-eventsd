
package main

import (
	"fmt"
	"log"
	"os"
	"io/ioutil"

	"gopkg.in/yaml.v1"
)

func loadYaml(filename string) (*EventsFile, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	eventsFile := &EventsFile{}
	err = yaml.Unmarshal(data, eventsFile)
	if err != nil {
		return nil, err
	}
	return eventsFile, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s filename.yml\n", os.Args[0])
		os.Exit(1)
	}
	filename := os.Args[1]
	eventsFile, err := loadYaml(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot load \"%s\": %s\n", filename, err)
		os.Exit(1)
	}
}
