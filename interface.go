// Package decomposer supports decimal marshalling and unmarshalling
// to and from components. Supports both fixed and arbitrary precision decimals.
package decomposer

// Decomposer composes or decomposes a decimal value to and from individual parts.
// There are four separate parts: a boolean negative flag, a form byte with three possible states
// (finite=0, infinite=1, NaN=2),  a base-2 little-endian integer
// coefficient (also known as a significand) as a []byte, and an int32 exponent.
// These are composed into a final value as "decimal = (neg) (form=finite) coefficient * 10 ^ exponent".
// A zero length coefficient is a zero value.
// If the form is not finite, the coefficient and scale should be zero and should be ignored.
// The negative parameter may be set to true for any form, although implementations are not required
// to respect the negative parameter in the non-finite form.
//
// Implementations may choose to signal a negative zero or negative NaN, but implementations
// that do not support these may also ignore the negative zero or negative NaN without error.
// If an implementation does not support Infinity it may be converted into a NaN without error.
// If a value is set that is larger then what is supported by an implementation is attempted to
// be set, an error must be returned.
// Implementations must return an error if a NaN or Infinity is attempted to be set while neither
// are supported.
type Decomposer interface {
	// Decompose returns the internal decimal state into parts.
	// If the provided buf has sufficient capacity, buf may be returned as the coefficient with
	// the value set and length set as appropriate.
	Decompose(buf []byte) (form byte, negative bool, coefficient []byte, exponent int32)

	// Compose sets the internal decimal value from parts. If the value cannot be
	// represented then an error should be returned.
	Compose(form byte, negative bool, coefficient []byte, exponent int32) error
}
