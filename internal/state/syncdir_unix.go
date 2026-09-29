//go:build !windows

package state

import "os"

// syncDir flushes a directory entry, so a rename into it survives a power cut.
func syncDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer func() { _ = d.Close() }()
	return d.Sync()
}
