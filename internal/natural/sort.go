package natural

import (
	"unicode"
	"unicode/utf8"
)

func Less(a, b string) bool {
	ai, bi := 0, 0
	for ai < len(a) && bi < len(b) {
		ar, aw := utf8.DecodeRuneInString(a[ai:])
		br, bw := utf8.DecodeRuneInString(b[bi:])
		ar, br = unicode.ToLower(ar), unicode.ToLower(br)

		ai += aw
		bi += bw

		adigit, bdigit := unicode.IsDigit(ar), unicode.IsDigit(br)

		// handle alphabet
		if !adigit && !bdigit {
			if ar != br {
				return ar < br
			}
			continue
		}

		if adigit != bdigit {
			return ar < br
		}

		anumber, bnumber := parsenum(a[ai-aw:]), parsenum(b[bi-bw:])
		if len(anumber) != len(bnumber) {
			return len(anumber) < len(bnumber)
		}
		if anumber != bnumber {
			return anumber < bnumber
		}

		ai += -aw + len(anumber)
		bi += -bw + len(bnumber)
	}

	if len(a) == len(b) {
		return false
	}
	return len(a) < len(b)
}

func parsenum(a string) string {
	for i, r := range a {
		if !unicode.IsDigit(r) {
			return a[:i]
		}
	}
	return a
}
