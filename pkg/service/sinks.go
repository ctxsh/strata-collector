package service

import (
	"reflect"

	"ctx.sh/strata-collector/pkg/apis/strata.ctx.sh/v1beta1"
	"ctx.sh/strata-collector/pkg/sink"
	"ctx.sh/strata-collector/pkg/sink/nats"
	"ctx.sh/strata-collector/pkg/sink/stdout"
)

func FromObject(obj v1beta1.CollectorOutput) sink.Sink {
	var s any
	// Iterate over the possible output configurations and choose the first non-nil
	// config.  The validation step should ensure that only one output is configured.
	output := reflect.ValueOf(obj)
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
