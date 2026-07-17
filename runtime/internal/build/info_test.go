package build

import "testing"

func TestCurrent(t *testing.T) {
	got := Current()

	if got.Version != Version {
		t.Errorf("Version = %q, want %q", got.Version, Version)
	}
	if got.Commit != Commit {
		t.Errorf("Commit = %q, want %q", got.Commit, Commit)
	}
	if got.Date != Date {
		t.Errorf("Date = %q, want %q", got.Date, Date)
	}
}
