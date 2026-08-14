package config

import (
	"time"

	"gopkg.in/yaml.v3"
)

// Duration is a time.Duration written in YAML as a Go duration string, e.g.
// "90s" or "500ms".
//
// It deliberately keeps the raw text and defers parsing to validation. YAML
// unmarshal errors are reported by line number only, whereas validation knows
// the field path — so "charger.heartbeat_timeout: ..." beats "line 5: ...".
type Duration struct {
	raw string
	dur time.Duration
}

// NewDuration returns a Duration with the given value, for tests and defaults.
func NewDuration(d time.Duration) Duration {
	return Duration{raw: d.String(), dur: d}
}

// Duration returns the parsed value. It is only meaningful after the config
// has passed validation.
func (d Duration) Duration() time.Duration { return d.dur }

// String returns the value as written by the operator.
func (d Duration) String() string { return d.raw }

// UnmarshalYAML captures the scalar as written. A non-scalar node leaves raw
// empty, which validation reports as a malformed duration.
func (d *Duration) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		d.raw = node.Value
	}
	return nil
}

// MarshalYAML writes the duration back out in the form it was given.
func (d Duration) MarshalYAML() (any, error) { return d.raw, nil }

func (d Duration) unset() bool { return d.raw == "" }

func (d *Duration) setDefault(v time.Duration) {
	d.raw = v.String()
	d.dur = v
}

// parse resolves the raw text. Callers report the returned error against the
// field's own path.
func (d *Duration) parse() error {
	v, err := time.ParseDuration(d.raw)
	if err != nil {
		return err
	}
	d.dur = v
	return nil
}
