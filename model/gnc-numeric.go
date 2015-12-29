package model

import (
	"fmt"
	"strconv"
	"strings"
)

// GncNumeric is the Gnucash value
type GncNumeric struct {
	num int
	den int
}

var (
	// GncNumericZero represents numeric 0
	GncNumericZero = NewGncNumeric(0, 100)

	// GncNumericOne represents numeric 1
	GncNumericOne = NewGncNumeric(100, 100)
)

// NewGncNumeric creates a new GncNumeric with numerator a and denominator b.
func NewGncNumeric(a, b int) *GncNumeric {
	if b == 0 {
		panic(fmt.Errorf("Division by 0: '%d/0'", a))
	}
	if b < 0 {
		a, b = -a, -b
	}
	return &GncNumeric{num: a, den: b}
}

//func (z *Rat) SetString(s string) (*Rat, bool)

// NewGncNumericFromString creates a new GncNumeric from string
func NewGncNumericFromString(v string) (*GncNumeric, error) {
	n := GncNumeric{}
	err := n.FromString(v)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

// FromString initialize GncNumeric from string
func (n *GncNumeric) FromString(v string) error {

	idx := strings.IndexByte(v, '/')
	if idx < 0 {

		num1, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		n.num = num1
		n.den = 1
	} else {
		num1, err := strconv.Atoi(v[0:idx])
		if err != nil {
			return err
		}
		den1, err := strconv.Atoi(v[idx+1:])
		if err != nil {
			return err
		}
		if den1 == 0 {
			return fmt.Errorf("Division by 0: '%s'", v)
		}
		if den1 < 0 {
			num1, den1 = -num1, -den1
		}
		n.num = num1
		n.den = den1
	}
	return nil
}

/*
// UnmarshalXML method of xml.Unmarshaler interface
func (n *GncNumeric) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	return n.FromString(v)
}
*/

// String returns a string representation of GncNumeric
func (n *GncNumeric) String() string {
	// does not check for zero denominator
	return fmt.Sprintf("%.02f", float64(n.num)/float64(n.den))
	//return fmt.Sprintf("%d/%d", n.num, n.den)
}

// Sign returns:
//
//	-1 if n <  0
//	 0 if n == 0
//	+1 if n >  0
//
func (n *GncNumeric) Sign() int {
	// assert(n.den > 0)
	if n.num > 0 {
		return 1
	}
	if n.num == 0 {
		return 0
	}
	return -1
}

// Neg sets n to -n and returns n.
func (n *GncNumeric) Neg() *GncNumeric {
	n.num = -n.num
	return n
}

// Add function: n.Add(x) -> n += x and returns n
func (n *GncNumeric) Add(x *GncNumeric) *GncNumeric {
	if n.den == x.den {
		n.num += x.num
		return n
	}
	g := LCM(n.den, x.den)
	n.num = n.num*(g/n.den) + x.num*(g/x.den)
	n.den = g
	return n
}

// Sub function: n.Sub(x) -> n -= x and returns n
func (n *GncNumeric) Sub(x *GncNumeric) *GncNumeric {
	if n.den == x.den {
		n.num -= x.num
		return n
	}
	g := LCM(n.den, x.den)
	n.num = n.num*(g/n.den) - x.num*(g/x.den)
	n.den = g
	return n
}

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// GCD returns the greatest common divisor of a and b.
func GCD(a, b int) int {
	a = Abs(a)
	b = Abs(a)
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM returns the least common multiple of a and b.
func LCM(a, b int) int {
	a = Abs(a)
	b = Abs(a)
	g := GCD(a, b)
	return (a / g) * b
}
