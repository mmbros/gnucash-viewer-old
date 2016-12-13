package numeric

import "testing"

func TestNew(t *testing.T) {
	const N numint = 3
	for j := -N; j <= N; j++ {
		expected := j
		if j < 0 {
			expected = -j
		}
		z := New(j, j)
		if z.num != expected {
			t.Errorf("New.num: expected %d, got %d", expected, z.num)
		}
		if z.den != expected {
			t.Errorf("New.den: expected %d, got %d", expected, z.den)
		}
	}
}
func TestZero(t *testing.T) {
	var testCases = []struct {
		num, den numint
		expected bool
	}{
		{0, 1, true},
		{0, 0, true},
		{1, 0, true},
		{1, 1, false},
	}

	for _, tc := range testCases {
		z := New(tc.num, tc.den)
		actual := z.Zero()
		if actual != tc.expected {
			t.Errorf("Zero: num=%d, den=%d, expected %v, got %v", tc.num, tc.den, tc.expected, actual)
		}
	}
}
func TestEqual(t *testing.T) {
	var testCases = []struct {
		a, b     Numeric
		expected bool
	}{
		{Numeric{1, 1}, Numeric{1, 1}, true},
		{Numeric{1, 1}, Numeric{5, 5}, false},
		{Numeric{0, 1}, Numeric{1, 0}, true},
		{New(-1, 2), New(1, -2), true},
	}

	for _, tc := range testCases {
		actual := tc.a.Equal(&tc.b)
		if actual != tc.expected {
			t.Errorf("Equal: %s == %s, expected %v, got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}

func TestString(t *testing.T) {
	var testCases = []struct {
		num, den numint
		str      string
	}{
		{1, 1, "1"},
		{10, 1, "10"},
		{-10, 1, "-10"},
		{10, -1, "-10"},
		{-10, -1, "10"},
		{5, 0, "0"},
		{0, 0, "0"},
		{-5, 0, "0"},
		{250, 100, "250/100"},
		{-250, 100, "-250/100"},
		{250, -100, "-250/100"},
		{-250, -100, "250/100"},
	}

	for _, tc := range testCases {
		z := New(tc.num, tc.den)
		actual := z.String()
		if actual != tc.str {
			t.Errorf("String: expected %q, got %q", tc.str, actual)
		}

	}
}

func TestFromString(t *testing.T) {
	var testCases = []struct {
		num, den numint
		str      string
	}{
		{1, 1, "1"},
		{10, 1, "10"},
		{-10, 1, "-10"},
		{5, 0, "0"},
		{0, 0, "0"},
		{-5, 0, "0"},
		{250, 100, "250/100"},
		{-250, 100, "250/-100"},
		{-250, 100, "-250/100"},
		{250, 100, "-250/-100"},
	}

	for _, tc := range testCases {
		actual, err := FromString(tc.str)
		if err != nil {
			t.Errorf("FromString(%q): unexpectedi error: %s", tc.str, err.Error())
			continue
		}
		expected := New(tc.num, tc.den)

		if !actual.Equal(&expected) {
			t.Errorf("FromString(%q): expected %q, got %q", tc.str, expected, actual)
		}

	}
}

func TestAdd(t *testing.T) {
	var testCases = []struct {
		a, b     Numeric
		expected Numeric
	}{
		{New(150, 100), New(250, 100), New(400, 100)},
		{New(1, 2), New(1, 3), New(5, 6)},
		{New(1, 2), New(5, 10), New(10, 10)},
	}

	for _, tc := range testCases {
		actual := Add(&tc.a, &tc.b)
		if !actual.Equal(&tc.expected) {
			t.Errorf("Add: %s + %s, expected %v, got %v", tc.a, tc.b, tc.expected, actual)
		}
	}
}
