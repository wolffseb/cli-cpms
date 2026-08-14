// Package buildinfo carries the version stamped into the binary at link time.
package buildinfo

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

// Set by the release build with -ldflags "-X .../internal/buildinfo.Version=...".
// A plain `go build` leaves these at their defaults and falls back to the VCS
// stamps the toolchain embeds automatically.
var (
	Version = "dev"
	Commit  = ""
	Date    = ""
)

// Info describes the running binary.
type Info struct {
	Version   string
	Commit    string
	Date      string
	GoVersion string
	Platform  string
}

// Get returns the build information, filling gaps from the embedded VCS stamps.
func Get() Info {
	info := Info{
		Version:   Version,
		Commit:    Commit,
		Date:      Date,
		GoVersion: runtime.Version(),
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
	}

	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				if info.Commit == "" {
					info.Commit = s.Value
				}
			case "vcs.time":
				if info.Date == "" {
					info.Date = s.Value
				}
			case "vcs.modified":
				if s.Value == "true" && info.Commit != "" {
					info.Commit += "-dirty"
				}
			}
		}
	}
	return info
}

// String renders a one-line summary, e.g. "cpms dev (abc1234) go1.24.7 linux/amd64".
func (i Info) String() string {
	s := "cpms " + i.Version
	if i.Commit != "" {
		short := i.Commit
		if len(short) > 12 {
			short = short[:12]
		}
		s += fmt.Sprintf(" (%s)", short)
	}
	if i.Date != "" {
		s += " built " + i.Date
	}
	return s + " " + i.GoVersion + " " + i.Platform
}
