package config

import (
	"errors"
	"testing"
)

func TestFlagError_FormatsAndUnwrapsError(t *testing.T) {
	t.Parallel()

	wrapped := errors.New("bad flag")
	err := &FlagError{Err: wrapped}

	if got, want := err.Error(), "parse configuration flags: bad flag"; got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
	if got := err.Unwrap(); !errors.Is(got, wrapped) {
		t.Fatalf("Unwrap() = %v, want %v", got, wrapped)
	}
}

func TestFileError_FormatsAndUnwrapsError(t *testing.T) {
	t.Parallel()

	wrapped := errors.New("permission denied")
	err := &FileError{Path: "/config.json", Operation: "open", Err: wrapped}

	if got, want := err.Error(), "open configuration file \"/config.json\": permission denied"; got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
	if got := err.Unwrap(); !errors.Is(got, wrapped) {
		t.Fatalf("Unwrap() = %v, want %v", got, wrapped)
	}
}

func TestMissingValueError_FormatsMissingField(t *testing.T) {
	t.Parallel()

	err := &MissingValueError{Field: "workspace"}

	if got, want := err.Error(), "required configuration value workspace is missing"; got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestInvalidValueError_FormatsAndUnwrapsError(t *testing.T) {
	t.Parallel()

	t.Run("without wrapped error", func(t *testing.T) {
		t.Parallel()

		err := &InvalidValueError{Field: "local_uid", Reason: "must be numeric"}

		if got, want := err.Error(), "configuration value local_uid is invalid: must be numeric"; got != want {
			t.Fatalf("Error() = %q, want %q", got, want)
		}
		if got := err.Unwrap(); got != nil {
			t.Fatalf("Unwrap() = %v, want nil", got)
		}
	})

	t.Run("with wrapped error", func(t *testing.T) {
		t.Parallel()

		wrapped := errors.New("strconv failure")
		err := &InvalidValueError{Field: "local_gid", Reason: "must be numeric", Err: wrapped}

		if got, want := err.Error(), "configuration value local_gid is invalid: must be numeric: strconv failure"; got != want {
			t.Fatalf("Error() = %q, want %q", got, want)
		}
		if got := err.Unwrap(); !errors.Is(got, wrapped) {
			t.Fatalf("Unwrap() = %v, want %v", got, wrapped)
		}
	})
}
