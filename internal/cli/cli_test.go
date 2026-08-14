package cli_test

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/wolffseb/cli-cpms/internal/cli"
)

// run executes the command tree with the given args and returns what it wrote
// to stdout and stderr along with the error the top level would see.
func run(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()

	cmd := cli.NewRootCommand()
	var out, errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	cmd.SetArgs(args)

	err = cmd.Execute()
	return out.String(), errOut.String(), err
}

func TestConfigValidateAcceptsExampleConfig(t *testing.T) {
	t.Parallel()

	stdout, stderr, err := run(t, "config", "validate", "-c", "../../config.example.yaml")
	if err != nil {
		t.Fatalf("expected success, got %v (stderr: %s)", err, stderr)
	}

	// The summary has to state the things an operator must get right for the
	// station and the counterparty to reach us.
	wants := []string{
		"is valid",
		"ALP-HYC-001",
		"ws://<this-host>:9000/ocpp/ALP-HYC-001", // the station's CSMS URL
		"DE*FRY",
		"http://192.168.1.10:8080/ocpi",
		"2 EVSEs",
		"DE*FRY*E001*1",
	}
	for _, want := range wants {
		if !strings.Contains(stdout, want) {
			t.Errorf("summary does not mention %q\n--- summary ---\n%s", want, stdout)
		}
	}
}

func TestConfigValidateReportsEveryProblem(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := dir + "/broken.yaml"
	writeFile(t, path, `
charger:
  ip: 192.168.1.42
server:
  ocpp_bind: 0.0.0.0:9000
  ocpi_bind: 0.0.0.0:9000
auth:
  default_id_tag: "04A1B2C3D4"
ocpi:
  country_code: de
  party_id: FRY
  public_base_url: http://192.168.1.10:8080/ocpi
  token_a: "token-a-secret"
  token_c: "token-c-secret"
location:
  id: OFFICE-01
  address: Musterstrasse 1
  city: Berlin
  country: DEU
  coordinates:
    latitude: "52.520008"
    longitude: "13.404954"
  evses:
    - uid: E1
      evse_id: DE*FRY*E001*1
      ocpp_connector_id: 1
      connectors:
        - id: "1"
          standard: IEC_62196_T2_COMBO
          format: CABLE
          power_type: DC
          max_voltage: 920
          max_amperage: 500
`)

	stdout, stderr, err := run(t, "config", "validate", "-c", path)
	if err == nil {
		t.Fatal("expected an invalid config to fail")
	}
	if !errors.Is(err, cli.ErrSilent) {
		t.Errorf("expected ErrSilent so the top level does not double-report, got %v", err)
	}
	if stdout != "" {
		t.Errorf("problems belong on stderr, but stdout had: %q", stdout)
	}

	// Three independent mistakes: no charger id, clashing ports, lowercase
	// country code. All three must be named in one run.
	for _, want := range []string{"charger.id", "server.ocpi_bind", "ocpi.country_code", "3 problems found"} {
		if !strings.Contains(stderr, want) {
			t.Errorf("report does not mention %q\n--- report ---\n%s", want, stderr)
		}
	}
}

func TestConfigValidateMissingFile(t *testing.T) {
	t.Parallel()

	_, _, err := run(t, "config", "validate", "-c", "does-not-exist.yaml")
	if err == nil {
		t.Fatal("expected an error for a missing file")
	}
	if errors.Is(err, cli.ErrSilent) {
		t.Error("an I/O error should be reported by the top level, not silenced")
	}
	if !strings.Contains(err.Error(), "does-not-exist.yaml") {
		t.Errorf("error %q does not name the file", err)
	}
}

func TestVersion(t *testing.T) {
	t.Parallel()

	stdout, _, err := run(t, "version")
	if err != nil {
		t.Fatalf("version: %v", err)
	}
	if !strings.HasPrefix(stdout, "cpms ") {
		t.Errorf("version output %q should start with %q", stdout, "cpms ")
	}
	if !strings.Contains(stdout, "go1.") {
		t.Errorf("version output %q should name the Go version", stdout)
	}
}

func TestRootShowsHelp(t *testing.T) {
	t.Parallel()

	stdout, _, err := run(t)
	if err != nil {
		t.Fatalf("bare invocation should succeed, got %v", err)
	}
	for _, want := range []string{"cpms", "config", "version"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("help does not mention %q\n%s", want, stdout)
		}
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}
