// Copyright 2020 Paul Borman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package check is designed for testing, including table drive tests, to check
// the value of a returned error.  The return value of a check is either an
// empty string (they matched) or a string describing why the check failed.
// Errors can be checked against different types of values (see Error) depending
// on the use case.
//
// Example Usage:
//
//	// We just care that there was an error
//	if s := check.Error(err, true); s != "" {
//		t.Errorf("Calling myFunc: %s", s)
//	}
//
//	// We want err.Error() to contain "unknown type"
//	if s := check.Error(err, "unknown type"); s != "" {
//		t.Errorf("Calling myFunc: %s", s)
//	}
//
//	// Got must be io.EOF
//	if s := check.Error(err, io.EOF); s != "" {
//		t.Errorf("Calling myFunc: %s", s)
//	}
//
//	// Is a wrapped io.EOF
//	if s := check.IsError(err, io.EOF); s != "" {
//		t.Errorf("Calling myFunc: %s", s)
//	}
//
// When a string is passed, check.Error normally just checks to see that error
// message contains the string, case sensitive.  Strings cast to Case are
// checked case insensitive while Equal and CaseEqual require the entire error
// message to be matched either case sensitive or insensitive respectively.
package check

import (
	"errors"
	"fmt"
	"strings"
)

// Type Equal is a string that an error must match exactly.
type Equal string

// Type Case is a string that must be contained in the error case insensitive.
type Case string

// Type CaseEqual is a string that error must case insensitive match exactly.
type CaseEqual string

// error formats

const (
	unexpected = "got unexpected error %q"
	expected   = "did not get expected error %q"
	wrong      = "got error %q, want %q"
)

var sprintf = fmt.Sprintf

// Error compares error got to interface want returning an empty string if they
// match or an error string if they are different.  The type of want determines
// how the check is made.
//
//	error:     got must be exactly want
//	bool:      check for existance of error
//	string:    check if got.Error() contains want
//	Case:      check if got.Error() contains want, case insensitive
//	Equal:     check if got.Error() is want
//	CaseEqual: check if got.Error() is want, case insensitive
func Error(got error, want interface{}) string {
	switch want := want.(type) {
	case bool:
		switch want {
		case (got != nil):
			return ""
		case true:
			return sprintf("did not get expected error")
		default:
			return sprintf(unexpected, got)
		}
	case Equal:
		switch {
		case got == nil && want == "":
			return ""
		case got == nil:
			return sprintf(expected, want)
		case want == "":
			return sprintf(unexpected, got)
		case got.Error() != string(want):
			return sprintf(wrong, got, want)
		default:
			return ""
		}
	case CaseEqual:
		switch {
		case got == nil && want == "":
			return ""
		case got == nil:
			return sprintf(expected, want)
		case want == "":
			return sprintf(unexpected, got)
		case strings.ToLower(got.Error()) != strings.ToLower(string(want)):
			return sprintf(wrong, got, want)
		default:
			return ""
		}
	case Case:
		switch {
		case got == nil && want == "":
			return ""
		case got == nil:
			return sprintf(expected, want)
		case want == "":
			return sprintf(unexpected, got)
		case !strings.Contains(strings.ToLower(got.Error()), strings.ToLower(string(want))):
			return sprintf(wrong, got, want)
		default:
			return ""
		}
	case string:
		switch {
		case got == nil && want == "":
			return ""
		case got == nil:
			return sprintf(expected, want)
		case want == "":
			return sprintf(unexpected, got)
		case !strings.Contains(got.Error(), want):
			return sprintf(wrong, got, want)
		default:
			return ""
		}
	case nil:
		// A nil interface appears to the type switch as type nil.
		// This means want can be any interface
		// type rather than only the error interface.  It also means
		// in the error case below we know want != nil.
		switch {
		case got == nil:
			return ""
		default:
			return sprintf(unexpected, got)
		}
	case error:
		switch {
		case got == nil:
			return sprintf(expected, want)
		case want != got:
			return sprintf(wrong, got, want)
		default:
			return ""
		}
	default:
		return sprintf("Check does not support type %T", want)
	}
}

// ErrorCase returns the empty string if got.Error() contains want, case
// insensitive, otherwise it returns a string indicating the error.
func ErrorCase(got error, want string) string {
	return Error(got, Case(want))
}

// ErrorCaseEqual returns the empty string if got.Error() matches want, case
// insensitive, otherwise it returns a string indicating the error.
func ErrorCaseEqual(got error, want string) string {
	return Error(got, CaseEqual(want))
}

// ErrorEqual returns the empty string if got.Error() exactly matches want
// otherwise it returns a string indicating the error.
func ErrorEqual(got error, want string) string {
	return Error(got, Equal(want))
}

// Is returns the empty string if want is is or is wrapped in got
// otherwise it returns a string indicating the error.
func IsError(got, want error) string {
	switch {
	case got == nil && want == nil:
		return ""
	case got == nil:
		return sprintf(expected, want)
	case want == nil:
		return sprintf(unexpected, got)
	case !errors.Is(got, want):
		return sprintf(wrong, got, want)
	default:
		return ""
	}
}
