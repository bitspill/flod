// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2018 The Flo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package flojson_test

import (
	"testing"

	"github.com/bitspill/flod/flojson"
)

// TestErrorCodeStringer tests the stringized output for the ErrorCode type.
func TestErrorCodeStringer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   flojson.ErrorCode
		want string
	}{
		{flojson.ErrDuplicateMethod, "ErrDuplicateMethod"},
		{flojson.ErrInvalidUsageFlags, "ErrInvalidUsageFlags"},
		{flojson.ErrInvalidType, "ErrInvalidType"},
		{flojson.ErrEmbeddedType, "ErrEmbeddedType"},
		{flojson.ErrUnexportedField, "ErrUnexportedField"},
		{flojson.ErrUnsupportedFieldType, "ErrUnsupportedFieldType"},
		{flojson.ErrNonOptionalField, "ErrNonOptionalField"},
		{flojson.ErrNonOptionalDefault, "ErrNonOptionalDefault"},
		{flojson.ErrMismatchedDefault, "ErrMismatchedDefault"},
		{flojson.ErrUnregisteredMethod, "ErrUnregisteredMethod"},
		{flojson.ErrNumParams, "ErrNumParams"},
		{flojson.ErrMissingDescription, "ErrMissingDescription"},
		{0xffff, "Unknown ErrorCode (65535)"},
	}

	// Detect additional error codes that don't have the stringer added.
	if len(tests)-1 != int(flojson.TstNumErrorCodes) {
		t.Errorf("It appears an error code was added without adding an " +
			"associated stringer test")
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

// TestError tests the error output for the Error type.
func TestError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   flojson.Error
		want string
	}{
		{
			flojson.Error{Description: "some error"},
			"some error",
		},
		{
			flojson.Error{Description: "human-readable error"},
			"human-readable error",
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		result := test.in.Error()
		if result != test.want {
			t.Errorf("Error #%d\n got: %s want: %s", i, result,
				test.want)
			continue
		}
	}
}
