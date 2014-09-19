
package main

type EventsFile struct {
	Cluster map[string]string
	Events []Event
}

type Event struct {
	Type string
	Command string
	Source string
	Contains string
}
