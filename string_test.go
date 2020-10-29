// Copyright 2020 The Godror Authors
//
// SPDX-License-Identifier: Apache-2.0

package decomposer_test

import (
	"fmt"
	"testing"

	"github.com/golang-sql/decomposer"
)

func TestNumberDeCompose(t *testing.T) {
	DecimalTest(t, func(s string) decomposer.Decimal { n := decomposer.NumberAsString(s); return &n })
}

func DecimalTest(t *testing.T, fromString func(string) decomposer.Decimal) {
	p := make([]byte, 38)
	for i, s := range []string{
		"0",
		"1",
		"-2",
		"3.14",
		"-3.14",
		"1000",
		"3.456789",
		"0.01",
		"-0.09",
		"-0.89",
		"0.0000000001",
		"12345678901234567890123456789012345678",
	} {
		n := fromString(s)

		form, negative, coefficient, exponent := n.Decompose(p[:0])
		if want := s[0] == '-'; want != negative {
			t.Errorf("%d. Decompose(%q) got negative=%t, wanted %t", i, s, negative, want)
		}
		if err := n.Compose(form, negative, coefficient, exponent); err != nil {
			t.Errorf("%d. cannot compose %c/%t/% x/%d from %q", i, form, negative, coefficient, exponent, s)
		}
		if got := fmt.Sprintf("%v", n); got != s {
			t.Errorf("%d. got %q wanted %q", i, got, s)
		}
	}
}
