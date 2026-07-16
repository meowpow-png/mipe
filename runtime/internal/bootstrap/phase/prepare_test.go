package phase

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/meowpow-png/mipe/runtime/internal/config"
	"go.uber.org/zap"
)

func testConfig() config.Config {
	return config.Config{
		AgentName:   "test-agent",
		UserHome:    "/home/user",
		AgentHome:   "/agent/home",
		RuntimeHome: "/runtime",
		Workspace:   "/workspace",
		LocalUID:    "1000",
		LocalGID:    "1001",
		Command:     []string{"bash"},
	}
}

func TestParseOwnership_ParsesNumericUidAndGid(t *testing.T) {
	t.Parallel()

	uid, gid, err := parseOwnership(testConfig())
	if err != nil {
		t.Fatalf("parseOwnership() error = %v", err)
	}
	if uid != 1000 || gid != 1001 {
		t.Fatalf("uid/gid = %d/%d, want 1000/1001", uid, gid)
	}
}

func TestParseOwnership_ReturnsErrorForInvalidUidOrGid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		edit func(*config.Config)
		want string
	}{
		{name: "uid", edit: func(cfg *config.Config) { cfg.LocalUID = "abc" }, want: "local_uid must be a numeric user id"},
		{name: "gid", edit: func(cfg *config.Config) { cfg.LocalGID = "abc" }, want: "local_gid must be a numeric group id"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := testConfig()
			tt.edit(&cfg)

			_, _, err := parseOwnership(cfg)
			if err == nil {
				t.Fatal("parseOwnership() error = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %q, want to contain %q", err, tt.want)
			}
		})
	}
}

func TestPrepare_CreatesCopiesAndChownsInOrder(t *testing.T) {
	restore := replacePrepareSeams(t)
	defer restore()

	var calls []string
	createDirectory = func(path string) error {
		calls = append(calls, "create:"+path)
		return nil
	}
	copyContents = func(source string, destination string) error {
		calls = append(calls, "copy:"+source+"->"+destination)
		return nil
	}
	chownRecursive = func(path string, uid int, gid int) error {
		calls = append(calls, "chown:"+path)
		if uid != 1000 || gid != 1001 {
			t.Fatalf("uid/gid = %d/%d, want 1000/1001", uid, gid)
		}
		return nil
	}
	if err := Prepare(testConfig(), zap.NewNop()); err != nil {
		t.Fatalf("Prepare() error = %v", err)
	}
	want := []string{
		"create:/agent/home",
		"copy:/runtime/config->/agent/home",
		"chown:/home/user",
	}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("calls = %#v, want %#v", calls, want)
	}
}

func TestPrepare_SkipsAgentHomeWhenUnset(t *testing.T) {
	restore := replacePrepareSeams(t)
	defer restore()

	cfg := testConfig()
	cfg.AgentHome = ""

	var calls []string
	createDirectory = func(path string) error {
		calls = append(calls, "create:"+path)
		return nil
	}
	copyContents = func(source string, destination string) error {
		calls = append(calls, "copy:"+source+"->"+destination)
		return nil
	}
	chownRecursive = func(path string, uid int, gid int) error {
		calls = append(calls, "chown:"+path)
		return nil
	}
	if err := Prepare(cfg, zap.NewNop()); err != nil {
		t.Fatalf("Prepare() error = %v", err)
	}
	want := []string{"chown:/home/user"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("calls = %#v, want %#v", calls, want)
	}
}

func TestPrepare_StopsOnFirstFileOperationError(t *testing.T) {
	tests := []struct {
		name      string
		failCall  string
		wantCalls []string
	}{
		{name: "create", failCall: "create", wantCalls: []string{"create"}},
		{name: "copy", failCall: "copy", wantCalls: []string{"create", "copy"}},
		{name: "chown", failCall: "chown", wantCalls: []string{"create", "copy", "chown"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := replacePrepareSeams(t)
			defer restore()

			sentinel := errors.New("boom")
			var calls []string
			createDirectory = func(path string) error {
				calls = append(calls, "create")
				if tt.failCall == "create" {
					return sentinel
				}
				return nil
			}
			copyContents = func(source string, destination string) error {
				calls = append(calls, "copy")
				if tt.failCall == "copy" {
					return sentinel
				}
				return nil
			}
			chownRecursive = func(path string, uid int, gid int) error {
				calls = append(calls, "chown")
				if tt.failCall == "chown" {
					return sentinel
				}
				return nil
			}
			err := Prepare(testConfig(), zap.NewNop())
			if !errors.Is(err, sentinel) {
				t.Fatalf("Prepare() error = %v, want sentinel", err)
			}
			if !reflect.DeepEqual(calls, tt.wantCalls) {
				t.Fatalf("calls = %#v, want %#v", calls, tt.wantCalls)
			}
		})
	}
}

func replacePrepareSeams(t *testing.T) func() {
	t.Helper()

	originalCreate := createDirectory
	originalCopy := copyContents
	originalChown := chownRecursive

	return func() {
		createDirectory = originalCreate
		copyContents = originalCopy
		chownRecursive = originalChown
	}
}
