package service

import "sync/atomic"

type DiscoveryStats struct {
	ReadyCollectors   atomic.Int64
	TotalResources    atomic.Int64
	InFlightResources atomic.Int64
}

func NewDiscoveryStats() *DiscoveryStats {
	return &DiscoveryStats{}
}

func (s *DiscoveryStats) SetReadyCollectors(i int64) {
	s.ReadyCollectors.Add(i)
}

func (s *DiscoveryStats) SetTotalResources(i int64) {
	s.TotalResources.Add(i)
}

func (s *DiscoveryStats) SetInFlightResources(i int64) {
	s.InFlightResources.Add(i)
}

func (s *DiscoveryStats) Reset() {
	s.ReadyCollectors.Store(0)
	s.TotalResources.Store(0)
	s.InFlightResources.Store(0)
}
