package loadbench

import "testing"

func TestPercentile(t *testing.T) {
	if got := percentile([]float64{1, 2, 3, 4, 5}, 0.95); got != 5 {
		t.Fatalf("p95=%v, want 5", got)
	}
	if got := percentile(nil, 0.95); got != 0 {
		t.Fatalf("empty p95=%v, want 0", got)
	}
}
