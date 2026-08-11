package loadbench

import (
	"strings"
	"testing"
	"time"
)

func TestPercentile(t *testing.T) {
	if got := percentile([]float64{1, 2, 3, 4, 5}, 0.95); got != 5 {
		t.Fatalf("p95=%v, want 5", got)
	}
	if got := percentile(nil, 0.95); got != 0 {
		t.Fatalf("empty p95=%v, want 0", got)
	}
}

func TestValidateConfigRequiresPublicationSamples(t *testing.T) {
	config := Config{
		Targets:            []string{"http://node-a", "http://node-b"},
		Duration:           time.Second,
		WarmupDuration:     time.Second,
		Concurrency:        1,
		WarmupIterations:   1,
		MeasuredIterations: 2,
		Rate:               1,
		Burst:              1,
	}
	if err := validateConfig(config); err == nil || !strings.Contains(err.Error(), "three measured") {
		t.Fatalf("expected measured-iteration validation, got %v", err)
	}
}

func TestSummaryUsesMedianAsPublishedValue(t *testing.T) {
	samples := []float64{12, 10, 100}
	metric := newMetric("total_rps", "requests_per_second", "higher_is_better", samples, 0)
	if metric.Value != 12 || metric.Summary["median"] != 12 || metric.Summary["max"] != 100 {
		t.Fatalf("unexpected metric summary: %+v", metric)
	}
}
