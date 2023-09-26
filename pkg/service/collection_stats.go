package service

import "sync/atomic"

type CollectionStats struct {
	// RegisteredDiscoveries is the number of discovery services that are
	// registered to the collector.
	RegisteredDiscoveries atomic.Int64
	// InFlightResources is the number of queued resources that are ready to
	// be processed.
	InFlightResources atomic.Int64
	// TotalSent is the number of metrics that have been sent to the output
	// successfully.  This value is reset at the end of each update cycle.
	TotalSent atomic.Int64
	// TotalErrors is the umber of metrics that have failed to be sent to the
	// output.  This value is reset at the end of each update cycle.
	TotalErrors atomic.Int64
	// TotalFiltered is the number of metrics that have been filtered out by
	// the collector.  This value is reset at the end of each update cycle.
	TotalFiltered atomic.Int64
	// MetricsCollected is the number of metrics collected by the collector. This
	// value is reset at the end of each update cycle.
	MetricsCollected atomic.Int64
}

func NewCollectionStats() *CollectionStats {
	return &CollectionStats{}
}

func (s *CollectionStats) SetRegisteredDiscoveries(i int64) {
	s.RegisteredDiscoveries.Add(i)
}

func (s *CollectionStats) SetInFlightResources(i int64) {
	s.InFlightResources.Add(i)
}

func (s *CollectionStats) SetTotalSent(i int64) {
	s.TotalSent.Add(i)
}

func (s *CollectionStats) SetTotalErrors(i int64) {
	s.TotalErrors.Add(i)
}

func (s *CollectionStats) SetTotalFiltered(i int64) {
	s.TotalFiltered.Add(i)
}

func (s *CollectionStats) SetMetricsCollected(i int) {
	s.MetricsCollected.Add(int64(i))
}

func (s *CollectionStats) Reset() {
	s.TotalSent.Store(0)
	s.TotalErrors.Store(0)
	s.TotalFiltered.Store(0)
	s.MetricsCollected.Store(0)
}
