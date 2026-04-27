package cmd //nolint:testpackage

import (
	"testing"
)

func Test_isK6Module(t *testing.T) {
	t.Parallel()

	cases := []struct {
		module string
		want   bool
	}{
		{"go.k6.io/k6", true},
		{"go.k6.io/k6/v2", true},
		{"go.k6.io/k6/v3", true},
		{"go.k6.io/k6/v10", true},
		// v0 and v1 have no major version suffix in Go module paths
		{"go.k6.io/k6/v0", false},
		{"go.k6.io/k6/v1", false},
		// not a version suffix
		{"go.k6.io/k6/extra", false},
		{"go.k6.io/k6/v", false},
		{"go.k6.io/k6/v2.1", false},
		// unrelated modules
		{"github.com/grafana/xk6-faker", false},
		{"go.k6.io/k6extra", false},
	}

	for _, c := range cases {
		t.Run(c.module, func(t *testing.T) {
			t.Parallel()

			if got := isK6Module(c.module); got != c.want {
				t.Errorf("isK6Module(%q) = %v, want %v", c.module, got, c.want)
			}
		})
	}
}
