package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/citadel/citadel"
)

type EventHandler struct {
	eventsFile   *EventsFile
	handlerIndex map[string]int
}

func NewEventHandler(eventsFile *EventsFile) (*EventHandler, error) {
	handlerIndex := make(map[string]int)
	for i, ev := range eventsFile.Events {
		value := strings.Replace(ev.Type, " ", "", -1)
		types := strings.Split(value, ",")
		for _, typ := range types {
			if typ == "" {
				return nil, errors.New("event has an invalid field `type'")
			}
			handlerIndex[strings.ToLower(typ)] = i
		}
	}
	eventHandler := &EventHandler{eventsFile, handlerIndex}
	return eventHandler, nil
}

func (h *EventHandler) findEvent(eventType string) []*Event {
	events := []*Event{}
	for typ, i := range h.handlerIndex {
		if typ == eventType {
			events = append(events, &h.eventsFile.Events[i])
		}
	}
	return events
}

func buildEnviron(env []string, event *citadel.Event) []string {
	env = append(env, fmt.Sprintf("FROM_ENGINE=%s", event.Engine.ID))
	env = append(env, fmt.Sprintf("FROM_CONTAINER=%s", event.Container.Name))
	return env
}

func execCommand(command string, env []string) {
	cmd := exec.Command("sh", "-l", "-c", command)
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Cannot exec the event handler command \"%s\": %s", command, err)
	} else {
		log.Printf("Returned: %s", out)
	}
}

func (h *EventHandler) Handle(e *citadel.Event) error {
	evs := h.findEvent(e.Type)
	if len(evs) < 1 {
		log.Printf("Uncaught event: type %s from %s@%s on %s", e.Type,
			e.Engine.ID, e.Engine.Addr, e.Container.Name)
		return nil
	}
	env := []string{}
	for name, addr := range h.eventsFile.Cluster {
		envVar := fmt.Sprintf("ENGINE_%s=%s", strings.ToUpper(name), addr)
		env = append(env, envVar)
	}
	for _, ev := range evs {
		if ev.FromEngine != "" &&
			(strings.ToLower(e.Engine.ID) != strings.ToLower(ev.FromEngine)) {
			log.Printf("Expected event from engine \"%s\", got it from \"%s\". Ignoring.",
				ev.FromEngine, e.Engine.ID)
			return nil
		}
		if ev.FromContainer != "" &&
			(strings.ToLower(e.Container.Name) != strings.ToLower(ev.FromContainer)) {
			log.Printf("Expected event from container \"%s\", got it from \"%s\". Ignoring.",
				ev.FromContainer, e.Container.Name)
			return nil
		}
		if ev.Command != "" {
			env := buildEnviron(env, e)
			execCommand(ev.Command, env)
		}
		if ev.Log != "" {
			log.Printf("LOG: %s", ev.Log)
		}
	}
	// Return code is not used there
	return nil
}
