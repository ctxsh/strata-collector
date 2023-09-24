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

package filter

import (
	"ctx.sh/strata-collector/pkg/metric"
)

type FilterFunc func(*metric.Metric) bool

type Filter struct {
	f []FilterFunc
}

func New() *Filter {
	return &Filter{
		f: make([]FilterFunc, 0),
	}
}

func (f *Filter) Use(ff ...FilterFunc) {
	f.f = append(f.f, ff...)
}

func (f *Filter) Do(m *metric.Metric) bool {
	for _, fn := range f.f {
		if fn(m) {
			return true
		}
	}

	return false
}

func Exclude(f ...float64) FilterFunc {
	return func(m *metric.Metric) bool {
		var value float64

		switch m.Vtype {
		case metric.Gauge:
			value = m.Values["gauge"].(float64)
		case metric.Counter:
			value = m.Values["counter"].(float64)
		case metric.Untyped:
			return false
		case metric.Summary:
			return false
		case metric.Histogram:
			return false
		}

		for _, v := range f {
			// fmt.Printf("value: %f, v: %f\n", value, v)
			if value == v {
				// fmt.Println("returning true")
				return true
			}
		}

		return false
	}
}

func Clip(min, max float64, inclusive bool) FilterFunc {
	return func(m *metric.Metric) bool {
		var value float64

		switch m.Vtype {
		case metric.Gauge:
			value = m.Values["gauge"].(float64)
		case metric.Counter:
			value = m.Values["counter"].(float64)
		case metric.Untyped:
			return false
		case metric.Summary:
			return false
		case metric.Histogram:
			return false
		}

		// fmt.Printf("value: %f, min: %f max: %f\n", value, min, max)
		if inclusive && (value < min || value > max) {
			// fmt.Println("matched inclusive")
			return true
		} else if !inclusive && (value <= min || value >= max) {
			// fmt.Println("matched noninclusive")
			return true
		}

		return false
	}
}
