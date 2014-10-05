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

func (h *EventHandler) findEvent(eventType string) *Event {
	for typ, i := range h.handlerIndex {
		if typ == eventType {
			return &h.eventsFile.Events[i]
		}
	}
	return nil
}

func (h *EventHandler) buildEnviron(event *citadel.Event) []string {
	env := []string{}
	for name, addr := range h.eventsFile.Cluster {
		envVar := fmt.Sprintf("ENGINE_%s=%s", strings.ToUpper(name), addr)
		env = append(env, envVar)
	}
	env = append(env, fmt.Sprintf("FROM_ENGINE=%s", event.Engine.ID))
	env = append(env, fmt.Sprintf("FROM_CONTAINER=%s", event.Container.Name))
	return env
}

func execCommand(command string, env []string) {
	cmd := exec.Command("sh", "-l", "-c", command)
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Cannot exec the event handler command: %s", err)
	}
	log.Printf("Returned: %s", out)
}

func (h *EventHandler) Handle(e *citadel.Event) error {
	ev := h.findEvent(e.Type)
	if ev == nil {
		log.Printf("Uncaught event: type %s from %s@%s on %s", e.Type,
			e.Engine.ID, e.Engine.Addr, e.Container.Name)
		return nil
	}
	//TODO: build environment and executes the different handlers (log, command)
	log.Printf("GOT EVENT: %#v", e)
	if ev.Command != "" {
		env := h.buildEnviron(e)
		execCommand(ev.Command, env)
	}
	if ev.Log != "" {
		log.Printf("LOG: %s", ev.Log)
	}
	// Return code is not used there
	return nil
}
