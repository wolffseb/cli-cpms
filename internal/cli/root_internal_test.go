package cli

import (
	"path/filepath"
	"testing"
)

func TestResolvedStatePath(t *testing.T) {
	t.Parallel()

	cases := []struct {
		config, state, want string
	}{
		{config: "config.yaml", want: "state.json"},
		{config: filepath.Join("a", "config.yaml"), want: filepath.Join("a", "state.json")},
		{config: filepath.Join("a", "b", "cpms.yaml"), want: filepath.Join("a", "b", "state.json")},
		{config: filepath.Join("a", "config.yaml"), state: "mine.json", want: "mine.json"},
	}
	for _, c := range cases {
		opts := &options{configPath: c.config, statePath: c.state}
		if got := opts.resolvedStatePath(); got != c.want {
			t.Errorf("config %q, --state %q: got %q, want %q", c.config, c.state, got, c.want)
		}
	}
}
