
package main

import (
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
		log.Fatalf("Usage: %s filename.yml", os.Args[0])
	}
	filename := os.Args[1]
	eventsFile, err := loadYaml(filename)
	if err != nil {
		log.Fatalf("Cannot load \"%s\": %s", filename, err)
	}
	log.Printf("ok: %#v\n", eventsFile)
	cluster, err := NewCluster(eventsFile.Cluster)
	if err != nil {
		log.Fatalf("Cannot init the cluster: %s", err)
	}
	log.Printf("ok: %#v\n", cluster)
	eventHandler := &EventHandler{eventsFile.Events}
	if err := cluster.Events(eventHandler); err != nil {
		log.Fatalf("Cannot init the log handler: %s", err)
	}
}
