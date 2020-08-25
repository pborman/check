# check ![build status](https://travis-ci.org/pborman/check.svg?branch=master) [![GoDoc](https://godoc.org/github.com/pborman/check?status.svg)](http://godoc.org/github.com/pborman/check)

Package check is designed for testing, including table drive tests, to check
the value of a returned error.  The return value of a check is either an
empty string (they matched) or a string describing why the check failed.
Errors can be checked against different types of values depending
on the use case.
```
TYPE       CHECK MECHANISM
error:     got must be exactly want
bool:      check for existance of error
string:    check if got.Error() contains want
Case:      check if got.Error() contains want, case insensitive
Equal:     check if got.Error() is want
CaseEqual: check if got.Error() is want, case insensitive
```
Example Usage:
```
        // We just care that there was an error
        if s := check.Error(err, true); s != "" {
                t.Errorf("Calling myFunc: %s", s)
        }

        // We want err.Error() to contain "unknown type"
        if s := check.Error(err, "unknown type"); s != "" {
                t.Errorf("Calling myFunc: %s", s)
        }

        // Got must be io.EOF
        if s := check.Error(err, io.EOF); s != "" {
                t.Errorf("Calling myFunc: %s", s)
        }

        // Is a wrapped io.EOF
        if s := check.IsError(err, io.EOF); s != "" {
                t.Errorf("Calling myFunc: %s", s)
        }
```

Typically ```check.Error``` is used in table drive tests where the condition
to check against is an ```interface{}``` value:

```
for _, tt := range []struct {
        input string
        err   interface{}
}{
        {"expect an error", true},
        {"expect io.EOF", io.EOF},
        {`expect exact match on "This Error"`, check.Equal("This Error")},
        {`expect "an error" to be in the error`, "an error"},
} {
        if s := error.Check(myFunc(tt.input), tt.err); s != "" {
                t.Errorf("myFunc(%s): %s", tt.input, s)
        }
}
```

When a string is passed, check.Error normally just checks to see that error
message contains the string, case sensitive.  Strings cast to Case are
checked case insensitive while Equal and CaseEqual require the entire error
message to be matched either case sensitive or insensitive respectively.
