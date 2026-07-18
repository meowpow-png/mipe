package phase

import (
	"reflect"
	"testing"
)

func TestCommandAsUser_SkipsGosuForMatchingIdentity(t *testing.T) {
	restore := replaceCurrentIdentity(t, 1000, 1001)
	defer restore()

	name, args, err := commandAsUser("1000", "1001", "env", "HOME=/home/user")
	if err != nil {
		t.Fatalf("commandAsUser() error = %v", err)
	}
	if name != "env" || !reflect.DeepEqual(args, []string{"HOME=/home/user"}) {
		t.Fatalf("commandAsUser() = %q %#v, want env %#v", name, args, []string{"HOME=/home/user"})
	}
}

func TestCommandAsUser_UsesGosuForMismatchedIdentity(t *testing.T) {
	restore := replaceCurrentIdentity(t, 1000, 1000)
	defer restore()

	name, args, err := commandAsUser("1000", "1001", "env", "HOME=/home/user")
	if err != nil {
		t.Fatalf("commandAsUser() error = %v", err)
	}
	want := []string{"1000:1001", "env", "HOME=/home/user"}
	if name != "gosu" || !reflect.DeepEqual(args, want) {
		t.Fatalf("commandAsUser() = %q %#v, want gosu %#v", name, args, want)
	}
}

func replaceCurrentIdentity(t *testing.T, uid, gid int) func() {
	t.Helper()
	original := currentIdentity
	currentIdentity = func() (int, int) { return uid, gid }
	return func() { currentIdentity = original }
}
