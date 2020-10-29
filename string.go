// Copyright 2020 Tamás Gulácsi
//
// SPDX-License-Identifier: Apache-2.0

package decomposer

import (
	"math/big"
	"strings"
)

// NumberAsString represents a base 10 decimal number as a string.
type NumberAsString string

// Decompose the base 10 decimal number from its string representation.
func (N NumberAsString) Decompose(buf []byte) (form byte, negative bool, coefficient []byte, exponent int32) {
	s := string(N)
	mexp := strings.IndexByte(s, '.')
	if mexp >= 0 {
		s = s[:mexp] + s[mexp+1:]
		exponent = -int32(len(s) - mexp)
	}
	var i big.Int
	if _, ok := i.SetString(s, 10); !ok {
		return 2, false, nil, 0
	}
	switch i.Sign() {
	case 0:
		return 0, false, nil, 0
	case -1:
		negative = true
	}
	c := (i.BitLen() + 7) >> 3
	if c <= cap(buf) {
		buf = buf[:c]
	} else {
		buf = make([]byte, c)
	}
	return 0, negative, i.FillBytes(buf), exponent
}

// Compose a base 10 decimal number as a string.
func (N *NumberAsString) Compose(form byte, negative bool, coefficient []byte, exponent int32) error {
	// This implementation tries hard to avoid extra allocations.
	var i big.Int
	var start int
	length := 1 + len(coefficient)*3
	if negative {
		start = 1
	}
	if exponent < 0 {
		length++
	} else if exponent > 0 {
		length += int(exponent)
	}
	p := make([]byte, start, length)
	if start != 0 {
		p[0] = '-'
	}
	i.SetBytes(coefficient)
	p = i.Append(p, 10)
	for ; exponent > 0; exponent-- {
		p = append(p, '0')
	}
	if exponent < 0 {
		exp := int(-exponent)
		if plus := exp - len(p) + start + 1; plus > 0 {
			olen := len(p)
			p = append(p, make([]byte, plus+1)...)
			copy(p[start+1+plus:], p[start:olen])
			for i := 0; i < plus; i++ {
				p[start+1+i] = '0'
			}
			p[start] = '0'
			p[start+1] = '.'
		} else {
			p = append(p, p[len(p)-1])
			copy(p[len(p)-1-exp:len(p)-1], p[len(p)-1-exp-1:len(p)-1-1])
			p[len(p)-1-exp] = '.'
		}
	}
	*N = NumberAsString(string(p))
	return nil
}

func (N NumberAsString) String() string { return string(N) }

var _ = Decimal((*NumberAsString)(nil))
