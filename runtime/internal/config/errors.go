package config

import "fmt"

type FlagError struct {
	Err error
}

func (err *FlagError) Error() string {
	return fmt.Sprintf("parse configuration flags: %v", err.Err)
}

func (err *FlagError) Unwrap() error {
	return err.Err
}

type FileError struct {
	Path      string
	Operation string
	Err       error
}

func (err *FileError) Error() string {
	return fmt.Sprintf("%s configuration file %q: %v", err.Operation, err.Path, err.Err)
}

func (err *FileError) Unwrap() error {
	return err.Err
}

type MissingValueError struct {
	Field string
}

func (err *MissingValueError) Error() string {
	return fmt.Sprintf("required configuration value %s is missing", err.Field)
}

type InvalidValueError struct {
	Field  string
	Reason string
	Err    error
}

func (err *InvalidValueError) Error() string {
	if err.Err == nil {
		return fmt.Sprintf("configuration value %s is invalid: %s", err.Field, err.Reason)
	}
	return fmt.Sprintf("configuration value %s is invalid: %s: %v", err.Field, err.Reason, err.Err)
}

func (err *InvalidValueError) Unwrap() error {
	return err.Err
}
