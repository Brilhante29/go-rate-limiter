package redisstore

import "testing"

func TestAsInt64(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  int64
	}{
		{name: "integer", input: int64(7), want: 7},
		{name: "string", input: "8", want: 8},
		{name: "bytes", input: []byte("9"), want: 9},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := asInt64(test.input)
			if err != nil || got != test.want {
				t.Fatalf("got=%d err=%v, want=%d", got, err, test.want)
			}
		})
	}
	if _, err := asInt64(1.5); err == nil {
		t.Fatal("expected unsupported type error")
	}
}
