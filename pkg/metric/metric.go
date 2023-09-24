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
	Vtype     ValueType              `json:"vtype"`
}

func New(t time.Time, name string, tags map[string]string) *Metric {
	metric := &Metric{
		Name:      name,
		Tags:      tags,
		Timestamp: t,
		Vtype:     Unknown,
		Values:    make(map[string]interface{}),
	}

	return metric
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
