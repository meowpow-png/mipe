package build

import "testing"

func TestCurrent(t *testing.T) {
	got := Current()

	if got.Version != Version {
		t.Errorf("Version = %q, want %q", got.Version, Version)
	}
}
