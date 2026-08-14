package config

import (
	"fmt"
	"sort"
	"strings"
)

// FieldError is a single problem, anchored to the YAML path it came from.
type FieldError struct {
	// Path is the dotted YAML path, e.g. "location.evses[1].ocpp_connector_id".
	Path string
	// Msg says what is wrong with the value at Path.
	Msg string
}

func (e FieldError) Error() string { return e.Path + ": " + e.Msg }

// ValidationErrors is every problem found in one pass. Validation does not
// stop at the first failure: an operator fixing a config file would rather see
// all six mistakes at once than rerun six times.
type ValidationErrors []FieldError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "config: no problems"
	}
	var b strings.Builder
	if len(e) == 1 {
		b.WriteString("config: 1 problem:")
	} else {
		fmt.Fprintf(&b, "config: %d problems:", len(e))
	}
	for _, fe := range e {
		fmt.Fprintf(&b, "\n  - %s", fe.Error())
	}
	return b.String()
}

// Paths returns just the field paths, in report order. Tests use it to assert
// which fields were flagged without pinning exact wording.
func (e ValidationErrors) Paths() []string {
	paths := make([]string, 0, len(e))
	for _, fe := range e {
		paths = append(paths, fe.Path)
	}
	return paths
}

// sortByPath gives a stable report order regardless of check order.
func (e ValidationErrors) sortByPath() {
	sort.SliceStable(e, func(i, j int) bool { return e[i].Path < e[j].Path })
}

// validator accumulates field errors during a validation pass.
type validator struct {
	errs ValidationErrors
}

func (v *validator) add(path, format string, args ...any) {
	v.errs = append(v.errs, FieldError{Path: path, Msg: fmt.Sprintf(format, args...)})
}

// required flags an empty value and reports whether the caller should keep
// checking it.
func (v *validator) required(path, value string) bool {
	if strings.TrimSpace(value) == "" {
		v.add(path, "is required")
		return false
	}
	return true
}
