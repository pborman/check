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

package check

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	err1 := errors.New(`Err one`)
	err2 := errors.New(`Err two`)
	err1u := `ERR ONE`
	err2u := `ERR TWO`

	for _, tt := range []struct {
		name string
		got  error
		want interface{}
		out  string
	}{
		// All cases where we got no error and expected no error
		{
			name: `nil no-error`,
		}, {
			name: `bool no-error`,
			want: false,
		}, {
			name: `error no-error`,
			want: error(nil),
		}, {
			name: `string no-error`,
			want: ``,
		}, {
			name: `equal no-error`,
			want: Equal(``),
		}, {
			name: `case no-error`,
			want: Case(``),
		}, {
			name: `caseequal no-error`,
			want: CaseEqual(``),
		},

		// All cases where we got the error we expected
		{
			name: `bool expected`,
			got:  err1,
			want: true,
		}, {
			name: `error expected`,
			got:  err1,
			want: err1,
		}, {
			name: `string expected`,
			got:  err1,
			want: `one`,
		}, {
			name: `equal expected`,
			got:  err1,
			want: Equal(err1.Error()),
		}, {
			name: `case expected`,
			got:  err1,
			want: Case(`ONE`),
		}, {
			name: `caseequal expected`,
			got:  err1,
			want: CaseEqual(err1u),
		},

		// All cases where we got an unexpected error
		{
			name: `nil unexpected`,
			got:  err1,
			out:  sprintf(unexpected, err1),
		}, {
			name: `bool unexpected`,
			got:  err1,
			want: false,
			out:  sprintf(unexpected, err1),
		}, {
			name: `error unexpected`,
			got:  err1,
			want: error(nil),
			out:  sprintf(unexpected, err1),
		}, {
			name: `string unexpected`,
			got:  err1,
			want: "",
			out:  sprintf(unexpected, err1),
		}, {
			name: `equal unexpected`,
			got:  err1,
			want: Equal(""),
			out:  sprintf(unexpected, err1),
		}, {
			name: `case unexpected`,
			got:  err1,
			want: Case(""),
			out:  sprintf(unexpected, err1),
		}, {
			name: `caseequal unexpected`,
			got:  err1,
			want: CaseEqual(""),
			out:  sprintf(unexpected, err1),
		},

		// All cases where we didn't get an expected error
		{
			name: `bool expected`,
			want: true,
			out:  `did not get expected error`,
		}, {
			name: `error expected`,
			want: err1,
			out:  sprintf(expected, err1),
		}, {
			name: `string expected`,
			want: err1.Error(),
			out:  sprintf(expected, err1),
		}, {
			name: `case expected`,
			want: Case(err1.Error()),
			out:  sprintf(expected, err1),
		}, {
			name: `equal expected`,
			want: Equal(err1.Error()),
			out:  sprintf(expected, err1),
		}, {
			name: `case equal expected`,
			want: CaseEqual(err1.Error()),
			out:  sprintf(expected, err1),
		},

		// All cases we go the wrong error
		{
			name: "error wrong",
			got:  err1,
			want: err2,
			out:  sprintf(wrong, err1, err2),
		}, {
			name: "string wrong",
			got:  err1,
			want: err2.Error(),
			out:  sprintf(wrong, err1, err2),
		}, {
			name: "case wrong",
			got:  err1,
			want: Case(err2u),
			out:  sprintf(wrong, err1, err2u),
		}, {
			name: "equal wrong",
			got:  err1,
			want: Equal(err2.Error()),
			out:  sprintf(wrong, err1, err2),
		}, {
			name: "caseequal wrong",
			got:  err1,
			want: CaseEqual(err2u),
			out:  sprintf(wrong, err1, err2u),
		},
		{
			name: `bad type`,
			want: 1,
			out:  `Check does not support type int`,
		}, {
			name: `bad nil`,
			want: (*struct{})(nil),
			out:  `Check does not support type *struct {}`,
		},
	} {
		s := Error(tt.got, tt.want)
		if s != tt.out {
			t.Errorf(`%s: got %q, want %q`, tt.name, s, tt.out)
		}
		switch w := tt.want.(type) {
		case Case:
			s := ErrorCase(tt.got, string(w))
			if s != tt.out {
				t.Errorf(`case-%s: got %q, want %q`, tt.name, s, tt.out)
			}
		case Equal:
			s := ErrorEqual(tt.got, string(w))
			if s != tt.out {
				t.Errorf(`equal-%s: got %q, want %q`, tt.name, s, tt.out)
			}
		case CaseEqual:
			s := ErrorCaseEqual(tt.got, string(w))
			if s != tt.out {
				t.Errorf(`case-equal-%s: got %q, want %q`, tt.name, s, tt.out)
			}
		}
	}
}

type errtype struct {
	E string
}

func (e *errtype) Error() string {
	if e == nil {
		return ""
	}
	return e.E
}

func TestIsError(t *testing.T) {
	err1 := &errtype{E: "basic"}
	err2 := errors.New("error 2")
	wrap1 := fmt.Errorf("wrapped %w", err1)

	for _, tt := range []struct {
		name string
		got  error
		want error
		out  string
	}{
		{
			name: "all nil",
		}, {
			name: "unexpected",
			got:  err1,
			out:  sprintf(unexpected, err1),
		},
		{
			name: "expected",
			want: err1,
			out:  sprintf(expected, err1),
		},
		{
			name: "same",
			got:  err1,
			want: err1,
		},
		{
			name: "mismatch",
			got:  err1,
			want: err2,
			out:  sprintf(wrong, err1, err2),
		},
		{
			name: "mismatch2",
			got:  err2,
			want: err1,
			out:  sprintf(wrong, err2, err1),
		},
		{
			name: "wrapped",
			got:  wrap1,
			want: err1,
		},
		{
			name: "wrapped wrong",
			got:  wrap1,
			want: err2,
			out:  sprintf(wrong, wrap1, err2),
		},
	} {
		s := IsError(tt.got, tt.want)
		if s != tt.out {
			t.Errorf(`%s: got %q, want %q`, tt.name, s, tt.out)
		}
	}
}
