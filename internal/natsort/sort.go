// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package natsort provides natural string sorting.
package natsort

import (
	"strings"
)

// Compare implements natural string sorting, where numbers are compared numerically.
//
// For example, "file1.txt" < "file2.txt" < "file10.txt".
func Compare(a, b string) int {
	if a == b {
		return 0
	}

	for {
		prefix := textPrefixLength(a, b)
		a, b = a[prefix:], b[prefix:]
		if len(a) == 0 || len(b) == 0 {
			return notEqualCompare(len(a), len(b))
		}

		// Did we reach a component with numbers?
		if isdigit(a[0]) && isdigit(b[0]) { // if both are numbers
			// We have two numbers.
			ac, azeros := countDigits(a)
			bc, bzeros := countDigits(b)

			// If one has more non-zero digits then it's obviously larger.
			if ac-azeros != bc-bzeros {
				return notEqualCompare(ac-azeros, bc-bzeros)
			}

			// Comparing equal length digit-strings will give the
			// same result as converting them to numbers.
			r := strings.Compare(a[azeros:ac], b[bzeros:bc])
			if r != 0 {
				return r
			}

			// The one with fewer leading zeros is smaller.
			if azeros != bzeros && azeros+bzeros > 0 {
				return notEqualCompare(azeros, bzeros)
			}

			// If they are numbers, we just continue.
			a, b = a[ac:], b[bc:]
		} else {
			// We know we reached differing characters.
			if a[0] < b[0] {
				return -1
			} else {
				return 1
			}
		}
	}
}

// notEqualCompare compares a and b assuming they are not equal.
func notEqualCompare(a, b int) int {
	if a < b {
		return -1
	} else {
		return 1
	}
}

// Less implements natural string comparison, where numbers are compared numerically.
func Less(a, b string) bool {
	return Compare(a, b) < 0
}

// Greater implements natural string comparison, where numbers are compared numerically.
func Greater(a, b string) bool {
	return Compare(a, b) > 0
}

// textPrefixLength returns the length of the longest common prefix of a and b ignoring digits.
func textPrefixLength(a, b string) int {
	i := 0
	for {
		if i >= len(a) || i >= len(b) {
			return i
		}
		ca, cb := a[i], b[i]
		if ca != cb || isdigit(ca) {
			return i
		}
		i++
	}
}

// countDigits returns the number of prefix digits in s.
func countDigits(s string) (count, leadingZeros int) {
	foundNonZero := false
	for i, c := range []byte(s) {
		if !isdigit(c) {
			return i, leadingZeros
		}
		if !foundNonZero && c == '0' {
			leadingZeros++
		} else {
			foundNonZero = true
		}
	}
	return len(s), leadingZeros
}

// isdigit returns true if c is a digit.
func isdigit(c byte) bool {
	return '0' <= c && c <= '9'
}
