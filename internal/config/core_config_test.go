package config

import (
	"testing"
)

func TestConfig_IsLocal(t *testing.T) {
	tests := map[string]struct {
		config *CoreConfig
		want   bool
	}{
		"local": {
			config: &CoreConfig{Environment: envLocal},
			want:   true,
		},
		"development": {
			config: &CoreConfig{Environment: envDevelopment},
			want:   false,
		},
	}

	for name, te := range tests {
		te := te
		got := te.config.IsLocal()

		if got != te.want {
			t.Errorf("[%s] got: %v, want: %v", name, got, te.want)
		}
	}
}
