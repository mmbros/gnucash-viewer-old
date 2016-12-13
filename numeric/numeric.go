// Package numeric defines the Numeric type.
package numeric

import (
	"fmt"
	"strconv"
	"strings"
)

// Numeric type is used to represents a GnuCash numeric type.
// Numeric{} equals 0 numeric number.
type Numeric struct {
	// Numerator
	num numint

	// Denominator
	// if den == 0 then the Numeric is 0
	// den is always >= 0
	den numint
}

// String returns a string representation of Numeric.
func (z Numeric) String() string {
	switch z.den {
	case 0: // den == 0
		return "0"
	case 1: // den == 1
		return strconv.FormatInt(int64(z.num), 10)
	default:
		return fmt.Sprintf("%d/%d", z.num, z.den)
	}
}

// New creates a new numeric with numerator num and denominator den.
func New(num, den numint) Numeric {
	if den < 0 {
		num, den = -num, -den
	}
	return Numeric{num: num, den: den}
}

// FromString creates a new Numeric from string.
func FromString(v string) (Numeric, error) {
	var z Numeric

	idx := strings.IndexByte(v, '/')
	if idx < 0 {
		num1, err := atoi(v)
		if err != nil {
			return z, err
		}
		z.num = num1
		z.den = 1
	} else {
		num1, err := atoi(v[0:idx])
		if err != nil {
			return z, err
		}
		den1, err := atoi(v[idx+1:])
		if err != nil {
			return z, err
		}
		if den1 < 0 {
			num1, den1 = -num1, -den1
		}
		z.num = num1
		z.den = den1
	}
	return z, nil
}

// Set sets Numeric z to the value of Numeric x.
func (z *Numeric) Set(x *Numeric) {
	z.num, z.den = x.num, x.den
}

// Zero return true if Numeric is zero.
//
// NOTE: Numeric{1, 0} == Numeric{10, 0} == Numeric{0, 0} == Numeric{0, 1}.
func (z *Numeric) Zero() bool {
	return (z.num == 0) || (z.den == 0)
}

// Equal returns true if z == x.
//
// NOTE: Numeric{1, 1} != Numeric{10, 10}.
func (z *Numeric) Equal(x *Numeric) bool {
	if z.Zero() {
		return z.Zero()
	}
	return (z.num == x.num) && (z.den == x.den)
}

// Sign returns:
//
//	-1 if z <  0
//	 0 if z == 0
//	+1 if z >  0
//
func (z *Numeric) Sign() int {
	if z.den == 0 {
		return 0
	}
	// assert(den > 0)
	if z.num > 0 {
		return 1
	}
	if z.num == 0 {
		return 0
	}
	return -1
}

// NegEqual sets z to -z
func (z *Numeric) NegEqual() {
	z.num = -z.num
}

// AddEqual function: z.AddEqual(x) -> z += x
func (z *Numeric) AddEqual(x *Numeric) {

	if x.den == 0 {
		// z += 0
		return
	}
	if z.den == 0 {
		// 0 += x
		z.num = x.num
		z.den = x.den
		return
	}
	if z.den == x.den {
		z.num += x.num
		return
	}
	g := lcm(z.den, x.den)
	z.num = z.num*(g/z.den) + x.num*(g/x.den)
	z.den = g
}

// SubEqual function: z.SubEqual(x) -> z -= x
func (z *Numeric) SubEqual(x *Numeric) {
	y := Neg(x)
	z.AddEqual(&y)
}

// Add function returns x+y.
func Add(x *Numeric, y *Numeric) Numeric {
	z := *x
	z.AddEqual(y)
	return z
}

// Sub function returns x-y.
func Sub(x *Numeric, y *Numeric) Numeric {
	z := *x
	z.SubEqual(y)
	return z
}

// Neg function sets x to -x.
func Neg(x *Numeric) Numeric {
	return Numeric{num: -x.num, den: x.den}
}

// Float64 converts the Numeric to a float64 value.
func (z *Numeric) Float64() float64 {
	if z.num == 0 || z.den == 0 {
		return 0.0
	}
	return float64(z.num) / float64(z.den)
}
