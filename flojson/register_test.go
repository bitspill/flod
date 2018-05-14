// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2018 The Flo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package flojson_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/bitspill/flod/flojson"
)

// TestUsageFlagStringer tests the stringized output for the UsageFlag type.
func TestUsageFlagStringer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   flojson.UsageFlag
		want string
	}{
		{0, "0x0"},
		{flojson.UFWalletOnly, "UFWalletOnly"},
		{flojson.UFWebsocketOnly, "UFWebsocketOnly"},
		{flojson.UFNotification, "UFNotification"},
		{flojson.UFWalletOnly | flojson.UFWebsocketOnly,
			"UFWalletOnly|UFWebsocketOnly"},
		{flojson.UFWalletOnly | flojson.UFWebsocketOnly | (1 << 31),
			"UFWalletOnly|UFWebsocketOnly|0x80000000"},
	}

	// Detect additional usage flags that don't have the stringer added.
	numUsageFlags := 0
	highestUsageFlagBit := flojson.TstHighestUsageFlagBit
	for highestUsageFlagBit > 1 {
		numUsageFlags++
		highestUsageFlagBit >>= 1
	}
	if len(tests)-3 != numUsageFlags {
		t.Errorf("It appears a usage flag was added without adding " +
			"an associated stringer test")
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		result := test.in.String()
		if result != test.want {
			t.Errorf("String #%d\n got: %s want: %s", i, result,
				test.want)
			continue
		}
	}
}

// TestRegisterCmdErrors ensures the RegisterCmd function returns the expected
// error when provided with invalid types.
func TestRegisterCmdErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		method  string
		cmdFunc func() interface{}
		flags   flojson.UsageFlag
		err     flojson.Error
	}{
		{
			name:   "duplicate method",
			method: "getblock",
			cmdFunc: func() interface{} {
				return struct{}{}
			},
			err: flojson.Error{ErrorCode: flojson.ErrDuplicateMethod},
		},
		{
			name:   "invalid usage flags",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				return 0
			},
			flags: flojson.TstHighestUsageFlagBit,
			err:   flojson.Error{ErrorCode: flojson.ErrInvalidUsageFlags},
		},
		{
			name:   "invalid type",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				return 0
			},
			err: flojson.Error{ErrorCode: flojson.ErrInvalidType},
		},
		{
			name:   "invalid type 2",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				return &[]string{}
			},
			err: flojson.Error{ErrorCode: flojson.ErrInvalidType},
		},
		{
			name:   "embedded field",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct{ int }
				return (*test)(nil)
			},
			err: flojson.Error{ErrorCode: flojson.ErrEmbeddedType},
		},
		{
			name:   "unexported field",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct{ a int }
				return (*test)(nil)
			},
			err: flojson.Error{ErrorCode: flojson.ErrUnexportedField},
		},
		{
			name:   "unsupported field type 1",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct{ A **int }
				return (*test)(nil)
			},
			err: flojson.Error{ErrorCode: flojson.ErrUnsupportedFieldType},
		},
		{
			name:   "unsupported field type 2",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct{ A chan int }
				return (*test)(nil)
			},
			err: flojson.Error{ErrorCode: flojson.ErrUnsupportedFieldType},
		},
		{
			name:   "unsupported field type 3",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct{ A complex64 }
				return (*test)(nil)
			},
			err: flojson.Error{ErrorCode: flojson.ErrUnsupportedFieldType},
		},
		{
			name:   "unsupported field type 4",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct{ A complex128 }
				return (*test)(nil)
			},
			err: flojson.Error{ErrorCode: flojson.ErrUnsupportedFieldType},
		},
		{
			name:   "unsupported field type 5",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct{ A func() }
				return (*test)(nil)
			},
			err: flojson.Error{ErrorCode: flojson.ErrUnsupportedFieldType},
		},
		{
			name:   "unsupported field type 6",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct{ A interface{} }
				return (*test)(nil)
			},
			err: flojson.Error{ErrorCode: flojson.ErrUnsupportedFieldType},
		},
		{
			name:   "required after optional",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct {
					A *int
					B int
				}
				return (*test)(nil)
			},
			err: flojson.Error{ErrorCode: flojson.ErrNonOptionalField},
		},
		{
			name:   "non-optional with default",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct {
					A int `jsonrpcdefault:"1"`
				}
				return (*test)(nil)
			},
			err: flojson.Error{ErrorCode: flojson.ErrNonOptionalDefault},
		},
		{
			name:   "mismatched default",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct {
					A *int `jsonrpcdefault:"1.7"`
				}
				return (*test)(nil)
			},
			err: flojson.Error{ErrorCode: flojson.ErrMismatchedDefault},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		err := flojson.RegisterCmd(test.method, test.cmdFunc(),
			test.flags)
		if reflect.TypeOf(err) != reflect.TypeOf(test.err) {
			t.Errorf("Test #%d (%s) wrong error - got %T, "+
				"want %T", i, test.name, err, test.err)
			continue
		}
		gotErrorCode := err.(flojson.Error).ErrorCode
		if gotErrorCode != test.err.ErrorCode {
			t.Errorf("Test #%d (%s) mismatched error code - got "+
				"%v, want %v", i, test.name, gotErrorCode,
				test.err.ErrorCode)
			continue
		}
	}
}

// TestMustRegisterCmdPanic ensures the MustRegisterCmd function panics when
// used to register an invalid type.
func TestMustRegisterCmdPanic(t *testing.T) {
	t.Parallel()

	// Setup a defer to catch the expected panic to ensure it actually
	// paniced.
	defer func() {
		if err := recover(); err == nil {
			t.Error("MustRegisterCmd did not panic as expected")
		}
	}()

	// Intentionally try to register an invalid type to force a panic.
	flojson.MustRegisterCmd("panicme", 0, 0)
}

// TestRegisteredCmdMethods tests the RegisteredCmdMethods function ensure it
// works as expected.
func TestRegisteredCmdMethods(t *testing.T) {
	t.Parallel()

	// Ensure the registered methods are returned.
	methods := flojson.RegisteredCmdMethods()
	if len(methods) == 0 {
		t.Fatal("RegisteredCmdMethods: no methods")
	}

	// Ensure the returned methods are sorted.
	sortedMethods := make([]string, len(methods))
	copy(sortedMethods, methods)
	sort.Sort(sort.StringSlice(sortedMethods))
	if !reflect.DeepEqual(sortedMethods, methods) {
		t.Fatal("RegisteredCmdMethods: methods are not sorted")
	}
}
