// Copyright 2023 Rob Lyon <rob@ctxswitch.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metric

import (
	"time"
)

// TODO: move me out to a common package

// MetricsType represents the type of metric that has been
// collected.
type MetricsType string

const (
	Counter   MetricsType = "counter"
	Gauge     MetricsType = "gauge"
	Untyped   MetricsType = "untyped"
	Summary   MetricsType = "summary"
	Histogram MetricsType = "histogram"
	Unknown   MetricsType = "unknown"
)

// Metric is used to store a scraped metric from prometheus
type Metric struct {
	// Name of the metric.
	Name string `json:"name"`
	// Tags are a map of key/value pairs that are used to store data that
	// is inteded to be indexed in upstream storage.
	Tags map[string]string `json:"tags"`
	// Timestamp represents the time that the metric was scraped.
	Timestamp time.Time `json:"timestamp"`
	// Type represents the type of metric that was scraped.
	Type MetricsType `json:"type"`
	// Value represents the value of the metric that was scraped as a float64.
	Value float64 `json:"value"`
}

// New creates a new metric.
func New(t time.Time, name string, value float64, tags map[string]string) *Metric {
	metric := &Metric{
		Name:      name,
		Timestamp: t,
		Type:      Unknown,
		Value:     value,
		Tags:      tags,
	}

	return metric
}

// AddField adds a tag to the metric.
func (m *Metric) AddTag(k, v string) {
	m.Tags[k] = v
}

// SetType sets the type of the metric.
func (m *Metric) SetType(Type MetricsType) {
	m.Type = Type
}
