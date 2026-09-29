package cli_test

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestStatePathDefaultsToTheConfigDirectory(t *testing.T) {
	t.Parallel()

	// The test binary runs in internal/cli, so a state path resolved against
	// the working directory would not land here.
	dir := filepath.Join(t.TempDir(), "a")
	if err := os.Mkdir(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(dir, "config.yaml")
	copyFile(t, "../../config.example.yaml", config)

	stdout, stderr, err := run(t, "config", "validate", "-c", config)
	if err != nil {
		t.Fatalf("validate: %v (stderr: %s)", err, stderr)
	}
	if want := filepath.Join(dir, "state.json"); !strings.Contains(stdout, want) {
		t.Errorf("summary does not name the state file %s:\n%s", want, stdout)
	}
}

func TestStateFlagOverridesTheDefault(t *testing.T) {
	t.Parallel()

	want := filepath.Join(t.TempDir(), "elsewhere", "cpms-state.json")
	stdout, stderr, err := run(t, "config", "validate", "-c", "../../config.example.yaml", "--state", want)
	if err != nil {
		t.Fatalf("validate: %v (stderr: %s)", err, stderr)
	}
	if !strings.Contains(stdout, want) {
		t.Errorf("summary does not name the state file %s:\n%s", want, stdout)
	}
	if strings.Contains(stdout, string(filepath.Separator)+"state.json") {
		t.Errorf("summary still mentions the default state file:\n%s", stdout)
	}
}

func TestRunRefusesACorruptStateFile(t *testing.T) {
	t.Parallel()

	const addr = "127.0.0.1:19104"
	config := runConfig(t, 19104)
	statePath := filepath.Join(filepath.Dir(config), "state.json")
	const corrupt = `{"version":1,`
	writeFile(t, statePath, corrupt)

	stdout, stderr, err := run(t, "run", "-c", config)
	if err == nil {
		t.Fatal("expected run to fail on a corrupt state file")
	}
	if !strings.Contains(err.Error(), statePath) {
		t.Errorf("error %q does not name the state file", err)
	}
	if out := stdout + stderr; strings.Contains(out, "listening") {
		t.Errorf("run started the listener before checking state:\n%s", out)
	}
	if c, err := net.DialTimeout("tcp", addr, 200*time.Millisecond); err == nil {
		_ = c.Close()
		t.Errorf("something is listening on %s", addr)
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != corrupt {
		t.Errorf("corrupt state file was modified: %q", data)
	}
}

func TestRunWithoutAStateFileCreatesNone(t *testing.T) {
	t.Parallel()

	config := runConfig(t, 19105)
	dir := filepath.Dir(config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	out, done := startRun(t, ctx, "run", "-c", config)
	waitForOutput(t, out, "cpms listening")
	cancel()
	if err := <-done; err != nil {
		t.Fatalf("run returned %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Name() != "config.yaml" {
			t.Errorf("run left %s behind in the config directory", e.Name())
		}
	}
}

func copyFile(t *testing.T, from, to string) {
	t.Helper()
	data, err := os.ReadFile(from)
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, to, string(data))
}
