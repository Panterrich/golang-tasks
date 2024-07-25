//go:build !solution

package testequal

import (
	"bytes"
	"maps"
	"slices"
)

func checkEqual(expected, actual interface{}) bool {
	switch e := expected.(type) {
	case int:
		if a, ok := actual.(int); ok && e == a {
			return true
		}
	case int8:
		if a, ok := actual.(int8); ok && e == a {
			return true
		}
	case int16:
		if a, ok := actual.(int16); ok && e == a {
			return true
		}
	case int32:
		if a, ok := actual.(int32); ok && e == a {
			return true
		}
	case int64:
		if a, ok := actual.(int64); ok && e == a {
			return true
		}
	case uint:
		if a, ok := actual.(uint); ok && e == a {
			return true
		}
	case uint8:
		if a, ok := actual.(uint8); ok && e == a {
			return true
		}
	case uint16:
		if a, ok := actual.(uint16); ok && e == a {
			return true
		}
	case uint32:
		if a, ok := actual.(uint32); ok && e == a {
			return true
		}
	case uint64:
		if a, ok := actual.(uint64); ok && e == a {
			return true
		}
	case string:
		if a, ok := actual.(string); ok && e == a {
			return true
		}
	case map[string]string:
		if a, ok := actual.(map[string]string); ok {
			if e == nil && a == nil {
				return true
			}

			if e != nil && a != nil && maps.Equal(e, a) {
				return true
			}

			return false
		}
	case []int:
		if a, ok := actual.([]int); ok {
			if e == nil && a == nil {
				return true
			}

			if e != nil && a != nil && slices.Equal(e, a) {
				return true
			}

			return false
		}
	case []byte:
		if a, ok := actual.([]byte); ok {
			if e == nil && a == nil {
				return true
			}

			if e != nil && a != nil && bytes.Equal(e, a) {
				return true
			}

			return false
		}
	}

	return false
}

func printMsgAndArgs(t T, msgAndArgs ...interface{}) {
	t.Helper()

	if len(msgAndArgs) == 0 {
		t.Errorf("")
	} else if len(msgAndArgs) == 1 {
		t.Errorf(msgAndArgs[0].(string))
	} else {
		t.Errorf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
}

// AssertEqual checks that expected and actual are equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are equal.
func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	if checkEqual(expected, actual) {
		return true
	}

	t.Helper()
	printMsgAndArgs(t, msgAndArgs...)

	return false
}

// AssertNotEqual checks that expected and actual are not equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are not equal.
func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	if !checkEqual(expected, actual) {
		return true
	}

	t.Helper()
	printMsgAndArgs(t, msgAndArgs...)

	return false
}

// RequireEqual does the same as AssertEqual but fails caller test immediately.
func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	if checkEqual(expected, actual) {
		return
	}

	t.Helper()
	printMsgAndArgs(t, msgAndArgs...)

	t.FailNow()
}

// RequireNotEqual does the same as AssertNotEqual but fails caller test immediately.
func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	if !checkEqual(expected, actual) {
		return
	}

	t.Helper()
	printMsgAndArgs(t, msgAndArgs...)

	t.FailNow()
}
