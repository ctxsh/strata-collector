package metric

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

func FromPrometheusMetric(now time.Time, buf []byte) ([]*Metric, error) {
	var parser expfmt.TextParser
	var err error

	var metrics []*Metric

	buf = bytes.TrimPrefix(buf, []byte("\n"))
	buffer := bytes.NewBuffer(buf)
	reader := bufio.NewReader(buffer)

	metricFamilies, err := parser.TextToMetricFamilies(reader)
	if err != nil {
		return nil, err
	}

	for name, mf := range metricFamilies {
		for _, m := range mf.Metric {
			tags := ParseLabelPairs(m.GetLabel())
			switch mf.GetType() {
			case dto.MetricType_SUMMARY:
				for _, q := range m.GetSummary().Quantile {
					if v := q.GetValue(); !math.IsNaN(v) {
						p := New(now, name, v, tags)
						p.SetType(Summary)

						quantile := fmt.Sprint(q.GetQuantile())
						p.AddTag("quantile", quantile)

						metrics = append(metrics, p)
					}
				}
			case dto.MetricType_HISTOGRAM:
				for _, b := range m.GetHistogram().Bucket {
					v := float64(b.GetCumulativeCount())
					p := New(now, name, v, tags)
					p.SetType(Histogram)

					bucket := fmt.Sprint(b.GetUpperBound())
					p.AddTag("bucket", bucket)

					metrics = append(metrics, p)
				}
			case dto.MetricType_COUNTER:
				if v := m.GetCounter().GetValue(); !math.IsNaN(v) {
					p := New(now, name, v, tags)
					p.SetType(Counter)

					metrics = append(metrics, p)
				}
			case dto.MetricType_GAUGE:
				if v := m.GetGauge().GetValue(); !math.IsNaN(v) {
					p := New(now, name, v, tags)
					p.SetType(Gauge)

					metrics = append(metrics, p)
				}
			case dto.MetricType_UNTYPED:
				if v := m.GetUntyped().GetValue(); !math.IsNaN(v) {
					p := New(now, name, v, tags)
					p.SetType(Untyped)

					metrics = append(metrics, p)
				}
			default:
				continue
			}
		}
	}

	return metrics, nil
}

func ParseLabelPairs(pairs []*dto.LabelPair) map[string]string {
	tags := make(map[string]string)

	for _, pair := range pairs {
		name := pair.GetName()
		value := pair.GetValue()
		tags[name] = value
	}

	return tags
}
