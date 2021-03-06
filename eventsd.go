package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v1"
)

var ExitChan chan bool

type EventsFile struct {
	Cluster map[string]string
	Events  []Event
}

type Event struct {
	Type          string `yaml:"type,omitempty"`
	Command       string `yaml:"command,omitempty"`
	FromEngine    string `yaml:"from_engine,omitempty"`
	FromContainer string `yaml:"from_container,omitempty"`
	ImageContains string `yaml:"image_contains,omitempty"`
	Log           string `yaml:"log,omitempty"`
}

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
		log.Fatalf("Usage: %s filename.yml", os.Args[0])
	}
	filename := os.Args[1]
	eventsFile, err := loadYaml(filename)
	if err != nil {
		log.Fatalf("Cannot load \"%s\": %s", filename, err)
	}
	cluster, err := NewCluster(eventsFile.Cluster)
	if err != nil {
		log.Fatalf("Cannot init the cluster: %s", err)
	}
	ExitChan = make(chan bool, 1)
	eventHandler, err := NewEventHandler(eventsFile)
	if err != nil {
		log.Fatalf("Cannot create the event handler: %s", err)
	}
	if err := cluster.Events(eventHandler); err != nil {
		log.Fatalf("Cannot init the log handler: %s", err)
	}
	log.Println("Listening to events...")
	<-ExitChan
	log.Println("Stop.")
}
