package service

import (
	"encoding/json"
	"time"
)

// TODO: move me out to a common package

// ValueType represents the type of metric that has been
// collected.
type ValueType string

const (
	Counter   ValueType = "counter"
	Gauge     ValueType = "gauge"
	Untyped   ValueType = "untyped"
	Summary   ValueType = "summary"
	Histogram ValueType = "histogram"
	Unknown   ValueType = "unknown"
)

// Metric is used to store a scraped metric from prometheus
type Metric struct {
	Name      string                 `json:"name"`
	Values    map[string]interface{} `json:"values"`
	Tags      map[string]string      `json:"tags"`
	Timestamp time.Time              `json:"timestamp"`
	Prefix    string                 `json:"prefix"`
	Vtype     ValueType              `json:"vtype"`
}

func NewMetric(t time.Time, name string, tags map[string]string) *Metric {
	metric := &Metric{
		Name:      name,
		Tags:      tags,
		Timestamp: t,
		Vtype:     Unknown,
		Values:    make(map[string]interface{}),
	}

	return metric
}

func (m *Metric) WithPrefix(prefix string) *Metric {
	m.Prefix = prefix
	return m
}

func (m *Metric) AddTag(k, v string) {
	m.Tags[k] = v
}

func (m *Metric) SetType(vtype ValueType) {
	m.Vtype = vtype
}

func (m *Metric) AddValue(name string, value interface{}) {
	m.Values[name] = value
}

func (m *Metric) Bytes() []byte {
	// TODO: this will most likely be added to the encoders, however
	// at this point in time, we don't quite know what the final relationship
	// will be.
	data, _ := json.Marshal(m)
	return data
}
