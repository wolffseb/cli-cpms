package state

// syncDir is a no-op: Windows has no equivalent of fsync on a directory.
func syncDir(string) error { return nil }
