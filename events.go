
package main

import (
	"errors"
	"log"
	"strings"

	"github.com/citadel/citadel"
)

type EventHandler struct {
	Events []Event
	handlerIndex map[string]int
}

func NewEventHandler(events []Event) (*EventHandler, error) {
	handlerIndex := make(map[string]int)
	for i, ev := range events {
		value := strings.Replace(ev.Type, " ", "", -1)
		types := strings.Split(value, ",")
		for _, typ := range types {
			if typ == "" {
				return nil, errors.New("event has an invalid field `type'")
			}
			handlerIndex[strings.ToLower(typ)] = i
		}
	}
	eventHandler := &EventHandler{events, handlerIndex}
	return eventHandler, nil
}

func (h *EventHandler) findEvent(eventType string) *Event {
	for typ, i := range h.handlerIndex {
		if typ == eventType {
			return &h.Events[i]
		}
	}
	return nil
}

func (h *EventHandler) Handle(e *citadel.Event) error {
	ev := h.findEvent(e.Type)
	if ev == nil {
		log.Printf("Uncaught event: type %s", e.Type)
		return nil
	}
	//TODO: build environment and executes the different handlers (log, command)
	return nil
}
