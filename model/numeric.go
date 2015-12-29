package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// base type of Numeric
type numint int64

// Numeric type
type Numeric struct {
	num numint
	den numint
}

var (
	// Numeric0 represents numeric 0
	Numeric0 = Numeric{num: 0, den: 100}

	// Numeric1 represents numeric 1
	Numeric1 = Numeric{num: 100, den: 100}
)

// String returns a string representation of Numeric
func (z Numeric) String() string {
	if z.den == 0 {
		if z.num == 0 {
			return "NaN"
		}
		if z.num > 0 {
			return "+Inf"
		}
		return "-Inf"
	}
	// den != 0
	return fmt.Sprintf("%.02f", float64(z.num)/float64(z.den))
	//return fmt.Sprintf("%d/%d", z.num, z.den)
}

// NewNumeric creates a new numeric with numerator num and denominator den.
func NewNumeric(num, den numint) Numeric {
	/*
		if den == 0 {
			return panic("Division by 0")
		}
	*/
	if den < 0 {
		num, den = -num, -den
	}
	return Numeric{num: num, den: den}
}

// NewNumericFromString creates a new Numeric from string
func NewNumericFromString(v string) (Numeric, error) {
	var z Numeric
	err := z.FromString(v)
	return z, err
}

// FromString initialize Numeric from string
func (z *Numeric) FromString(v string) error {

	idx := strings.IndexByte(v, '/')
	if idx < 0 {
		num1, err := _atoi(v)
		if err != nil {
			return err
		}
		z.num = num1
		z.den = 1
	} else {
		num1, err := _atoi(v[0:idx])
		if err != nil {
			return err
		}
		den1, err := _atoi(v[idx+1:])
		if err != nil {
			return err
		}
		if den1 == 0 {
			return errors.New("Division by 0")
		}
		if den1 < 0 {
			num1, den1 = -num1, -den1
		}
		z.num = num1
		z.den = den1
	}
	return nil
}

// Sign returns:
//
//	-1 if z <  0
//	 0 if z == 0 (or z == -Inf, NaN, +Inf)
//	+1 if z >  0
//
func (z *Numeric) Sign() int {
	if z.den > 0 {
		if z.num > 0 {
			return 1
		}
		if z.num == 0 {
			return 0
		}
		return -1
	}

	if z.den < 0 {
		if z.num > 0 {
			return -1
		}
		if z.num == 0 {
			return 0
		}
		return 1
	}

	// den == 0
	return 0
}

// Neg sets n to -n
func (z *Numeric) Neg() {
	z.num = -z.num
}

// Add function: n.Add(x) -> n += x
func (z *Numeric) Add(x Numeric) {
	if z.den == x.den {
		z.num += x.num
		return
	}
	g := _lcm(z.den, x.den)
	z.num = z.num*(g/z.den) + x.num*(g/x.den)
	z.den = g
}

// Sub function: n.Sub(x) -> n -= x
func (z *Numeric) Sub(x Numeric) {
	if z.den == x.den {
		z.num -= x.num
		return
	}
	g := _lcm(z.den, x.den)
	z.num = z.num*(g/z.den) - x.num*(g/x.den)
	z.den = g
}

//*************************************************************
//*************************************************************
//*************************************************************

// _atoi converts a string to a numint
func _atoi(s string) (numint, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	return numint(i), err
}

// _abs returns the absolute value of i.
func _abs(i numint) numint {
	if i < 0 {
		return -i
	}
	return i
}

// _gcd returns the greatest common divisor of a and b.
func _gcd(a, b numint) numint {
	a = _abs(a)
	b = _abs(b)
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// _lcm returns the least common multiple of a and b.
func _lcm(a, b numint) numint {
	a = _abs(a)
	b = _abs(b)
	g := _gcd(a, b)
	l := (a / g) * b
	//	fmt.Printf("lcm(%d, %d) = %d\n", a, b, l)

	return l
}
