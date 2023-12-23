package main

import (
	"sync"
)

type Observation struct {
	mile      int
	timestamp int64
}

type RoadSpeedLimitEnforcer struct {
	vehicleObservations map[string]([]*Observation)
	road                int
	limit               int
	mu                  sync.Mutex
}

func (_ *RoadSpeedLimitEnforcer) New(road int, limit int) *RoadSpeedLimitEnforcer {
	return &RoadSpeedLimitEnforcer{
		vehicleObservations: make(map[string][]*Observation),
		road:                road,
		limit:               limit,
	}
}

func (r *RoadSpeedLimitEnforcer) AddObservation(plate string, mile int, timestamp int64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.vehicleObservations == nil {
		panic("vehicleObservations is empty!")
	}

	observations := r.vehicleObservations[plate]
	observations = append(observations, &Observation{mile: mile, timestamp: timestamp})
	r.vehicleObservations[plate] = observations
}
