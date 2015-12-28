package model

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

// GncNumeric is the Gnucash value
type GncNumeric struct {
	num int
	den int
}

// NewGncNumeric creates a new GncNumeric with numerator a and denominator b.
func NewGncNumeric(a, b int) *GncNumeric {
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
		n.num = num1
		n.den = den1
	}
	return nil
}

// UnmarshalXML method of xml.Unmarshaler interface
func (n *GncNumeric) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	return n.FromString(v)
}

// String returns a string representation of GncNumeric
func (n *GncNumeric) String() string {
	// does not check for zero denominator
	return fmt.Sprintf("%.02f", float64(n.num)/float64(n.den))
	//return fmt.Sprintf("%d/%d", n.num, n.den)
}

// Add function: n.Add(x) -> n += x
func (n *GncNumeric) Add(x *GncNumeric) {
	if n.den == x.den {
		n.num += x.num
		return
	}
	//panic("Not Implemented: GncNumeric.Add with different denominator")
	g := LCM(n.den, x.den)
	n.num = n.num*(g/n.den) + x.num*(g/x.den)
	n.den = g
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
