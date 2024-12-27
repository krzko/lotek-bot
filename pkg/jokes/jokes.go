package jokes

import (
	"math/rand"
	"time"
)

// Service handles joke operations
type Service struct {
	jokes []string
}

// NewService creates a new jokes service
func NewService() *Service {
	return &Service{
		jokes: []string{
			"Why do programmers prefer dark mode? Because light attracts bugs!",
			"Why do Kubernetes administrators never get lost? Because they always follow the NodePath!",
			"What's a SRE's favorite breakfast? YAML and eggs!",
			"What did the Prometheus query say to the time series? You've got a lot of explaining to do!",
			"Why did the OpenTelemetry collector go to therapy? It had too much emotional baggage to trace!",
		},
	}
}

// GetRandomJoke returns a random joke
func (s *Service) GetRandomJoke() string {
	rand.Seed(time.Now().UnixNano())
	return s.jokes[rand.Intn(len(s.jokes))]
}

// AddJoke adds a new joke to the collection
func (s *Service) AddJoke(joke string) {
	s.jokes = append(s.jokes, joke)
}
