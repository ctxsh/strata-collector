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
