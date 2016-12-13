package numeric

import "testing"

func TestGcd(t *testing.T) {
	var testCases = []struct {
		a, b, expected numint
	}{
		{7, 3, 1},
		{7, -3, 1},
		{-7, -3, 1},
		{-7, 3, 1},
		{15, 6, 3},
		{15, -6, 3},
		{-15, -6, 3},
		{-15, 6, 3},
		{2, 0, 2},
		{0, 2, 2},
		{-2, 0, 2},
		{0, -2, 2},
		{0, 0, 0},
		{65535, 500, 5},
	}

	for _, tc := range testCases {
		actual := gcd(tc.a, tc.b)
		if actual != tc.expected {
			t.Errorf("gcd: a=%d, b=%d, expected %d, got %d", tc.a, tc.b, tc.expected, actual)
		}
	}
}

func TestLcm(t *testing.T) {
	var testCases = []struct {
		a, b, expected numint
	}{
		{7, 3, 21},
		{7, -3, 21},
		{-7, -3, 21},
		{-7, 3, 21},
		{15, 6, 30},
		{15, -6, 30},
		{-15, -6, 30},
		{-15, 6, 30},
		{2, 0, 0},
		{0, 2, 0},
		{-2, 0, 0},
		{0, -2, 0},
		{65535, 500, 6553500},
		{0, 0, 0},
	}

	for _, tc := range testCases {
		actual := lcm(tc.a, tc.b)
		if actual != tc.expected {
			t.Errorf("lcm: a=%d, b=%d, expected %d, got %d", tc.a, tc.b, tc.expected, actual)
		}
	}
}
