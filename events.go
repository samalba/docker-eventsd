
package main

import (
	"log"

	"github.com/citadel/citadel"
)

type EventHandler struct {
	Events []Event
}

func (h *EventHandler) Handle(e *citadel.Event) error {
	log.Printf("type: %s image: %s container: %s\n",
		e.Type, e.Container.Image.Name, e.Container.ID)
	return nil
}
