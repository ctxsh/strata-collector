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

package mutation

import "ctx.sh/strata-collector/pkg/metric"

type MutateFunc func(*metric.Metric)

type Mutator struct {
	m []MutateFunc
}

func New() *Mutator {
	return &Mutator{
		m: make([]MutateFunc, 0),
	}
}

func (mut *Mutator) Use(mf ...MutateFunc) {
	mut.m = append(mut.m, mf...)
}

func (mut *Mutator) Do(m *metric.Metric) {
	for _, fn := range mut.m {
		fn(m)
	}
}

// func Clamp(min, max, minVal, maxVal float64) MutateFunc {
// 	return func(m *metric.Metric) {
// 		var value float64

// 		switch m.Vtype {
// 		case metric.Gauge:
// 			value = m.Values["gauge"].(float64)
// 		case metric.Counter:
// 			value = m.Values["counter"].(float64)
// 		case metric.Untyped:
// 			return true
// 		case metric.Summary:
// 			return true
// 		case metric.Histogram:
// 			return true
// 		}

// 		for

// 		if value < min || value > max {
// 			m.Values["value"] = minVal
// 		}

// 		return true
// 	}
// }
