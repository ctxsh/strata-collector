package service

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

type PrometheusScraper struct {
	Url    string
	Client http.Client
}

func NewPrometheusScraper(client http.Client, url string) *PrometheusScraper {
	return &PrometheusScraper{
		Url:    url,
		Client: client,
	}
}

func (p *PrometheusScraper) Get(tags map[string]string) ([]*Metric, error) {
	req, _ := http.NewRequest("GET", p.Url, nil)
	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	m, err := p.parse(time.Now(), buf, tags)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (p *PrometheusScraper) parse(now time.Time, buf []byte, tags map[string]string) ([]*Metric, error) {
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
			tags := parseLabels(m, tags)
			tags = parseLabelPairs(tags, m.GetLabel())
			p := NewMetric(now, name, tags)

			switch mf.GetType() {
			// Parse summary metrics
			case dto.MetricType_SUMMARY:
				for _, q := range m.GetSummary().Quantile {
					if v := q.GetValue(); !math.IsNaN(v) {
						p.AddValue(fmt.Sprint(q.GetQuantile()), v)
					}
				}
			// Parse histogram metrics
			case dto.MetricType_HISTOGRAM:
				p.SetType(Histogram)
				for _, b := range m.GetHistogram().Bucket {
					p.AddValue(fmt.Sprint(b.GetUpperBound()), float64(b.GetCumulativeCount()))
				}
			// Parse counter metrics
			case dto.MetricType_COUNTER:
				if v := m.GetCounter().GetValue(); !math.IsNaN(v) {
					p.SetType(Counter)
					p.AddValue("counter", v)
				}
			// Parse gauge metrics
			case dto.MetricType_GAUGE:
				if v := m.GetGauge().GetValue(); !math.IsNaN(v) {
					p.SetType(Gauge)
					p.AddValue("gauge", v)
				}
			// Parse untyped metrics
			case dto.MetricType_UNTYPED:
				if v := m.GetUntyped().GetValue(); !math.IsNaN(v) {
					p.SetType(Untyped)
					p.AddValue("value", v)
				}
			default:
				continue
			}

			metrics = append(metrics, p)
		}
	}

	return metrics, nil
}

func parseLabels(m *dto.Metric, tags map[string]string) map[string]string {
	result := map[string]string{}

	for key, value := range tags {
		result[key] = value
	}

	for _, pair := range m.Label {
		result[pair.GetName()] = pair.GetValue()
	}

	return result
}

func parseLabelPairs(tags map[string]string, pairs []*dto.LabelPair) map[string]string {
	for _, pair := range pairs {
		name := pair.GetName()
		value := pair.GetValue()
		tags[name] = value
	}
	return tags
}
