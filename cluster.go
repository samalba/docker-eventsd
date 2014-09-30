package main

import (
	"github.com/citadel/citadel"
	"github.com/citadel/citadel/cluster"
	"github.com/citadel/citadel/scheduler"
)

func NewCluster(engines map[string]string) (*cluster.Cluster, error) {
	c, err := cluster.New(scheduler.NewResourceManager())
	if err != nil {
		return nil, err
	}
	for name, url := range engines {
		e := &citadel.Engine{
			ID:     name,
			Addr:   url,
			Memory: 2048, // unused
			Cpus:   4,    // unused
			Labels: []string{name},
		}
		err = e.Connect(nil)
		if err != nil {
			return nil, err
		}
		err = c.AddEngine(e)
		if err != nil {
			c.Close()
			return nil, err
		}
	}
	return c, nil
}
