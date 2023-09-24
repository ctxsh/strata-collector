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

package service

import (
	"reflect"

	"ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	"ctx.sh/strata-collector/pkg/encoder"
	"ctx.sh/strata-collector/pkg/encoder/json"
	"ctx.sh/strata-collector/pkg/filter"
	"ctx.sh/strata-collector/pkg/output"
	"ctx.sh/strata-collector/pkg/output/nats"
	"ctx.sh/strata-collector/pkg/output/stdout"
)

func OutputFactory(obj *v1beta1.CollectorOutput) output.Output {
	var s any
	// Iterate over the possible output configurations and choose the first non-nil
	// config.  The validation step should ensure that only one output is configured.
	output := reflect.ValueOf(*obj)
	for i := 0; i < output.NumField(); i++ {
		if !output.Field(i).IsNil() && output.Field(i).Kind() != reflect.String {
			s = output.Field(i).Interface()
			break
		}
	}

	switch s.(type) {
	case *v1beta1.Nats:
		return nats.New()
	default:
		return stdout.New()
	}
}

func EncoderFactory(name string) encoder.Encoder {
	switch name {
	default:
		return json.New()
	}
}

func FilterFactory(obj *v1beta1.CollectorFilters) *filter.Filter {
	if obj == nil {
		return nil
	}

	f := filter.New()

	if obj.Exclude != nil {
		f.Use(filter.Exclude(obj.Exclude.Values...))
	}

	if obj.Clip != nil {
		f.Use(filter.Clip(*obj.Clip.Min, *obj.Clip.Max, *obj.Clip.Inclusive))
	}

	return f
}
