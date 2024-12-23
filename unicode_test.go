package unicodeid_test

import (
	"fmt"
	"math/rand/v2"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/nukilabs/unicodeid"
	"golang.org/x/text/unicode/rangetable"
)

func genString(pNonASCII int) string {
	ascii := func() rune { return rune(rand.IntN(0x7F)) }
	nonascii := func() rune { return rune(rand.IntN(0x10FFFF-0x80) + 0x80) }

	result := make([]rune, 500_000)
	for i := range result {
		if rand.IntN(100) < pNonASCII {
			result[i] = nonascii()
		} else {
			result[i] = ascii()
		}
	}
	return string(result)
}

var (
	unicodeRangeIdNeg      = rangetable.Merge(unicode.Pattern_Syntax, unicode.Pattern_White_Space)
	unicodeRangeIdStartPos = rangetable.Merge(unicode.Letter, unicode.Nl, unicode.Other_ID_Start)
	unicodeRangeIdContPos  = rangetable.Merge(unicodeRangeIdStartPos, unicode.Mn, unicode.Mc, unicode.Nd, unicode.Pc, unicode.Other_ID_Continue)
)

func isIdStartUnicode(r rune) bool {
	return unicode.Is(unicodeRangeIdStartPos, r) && !unicode.Is(unicodeRangeIdNeg, r)
}

func isIdPartUnicode(r rune) bool {
	return unicode.Is(unicodeRangeIdContPos, r) && !unicode.Is(unicodeRangeIdNeg, r) || r == '\u200C' || r == '\u200D'
}

func isIdentifierStart(chr rune) bool {
	return 'a' <= chr && chr <= 'z' || 'A' <= chr && chr <= 'Z' ||
		chr >= utf8.RuneSelf && isIdStartUnicode(chr)
}

func isIdentifierPart(chr rune) bool {
	return chr == '_' || 'a' <= chr && chr <= 'z' || 'A' <= chr && chr <= 'Z' ||
		'0' <= chr && chr <= '9' ||
		chr >= utf8.RuneSelf && isIdPartUnicode(chr)
}

func bench(b *testing.B, pNonASCII int) {
	str := genString(pNonASCII)
	b.Run(fmt.Sprintf("baseline-%d", pNonASCII), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, r := range str {
				_ = r
			}
		}
	})
	b.Run(fmt.Sprintf("unicodeid-%d", pNonASCII), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, r := range str {
				unicodeid.IsIDStart(r)
				unicodeid.IsIDContinue(r)
			}
		}
	})
	b.Run(fmt.Sprintf("unicode-%d", pNonASCII), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, r := range str {
				isIdentifierStart(r)
				isIdentifierPart(r)
			}
		}
	})
}

func Benchmark0NonASCII(b *testing.B) {
	bench(b, 0)
}

func Benchmark1NonASCII(b *testing.B) {
	bench(b, 1)
}

func Benchmark10NonASCII(b *testing.B) {
	bench(b, 10)
}

func Benchmark100NonASCII(b *testing.B) {
	bench(b, 100)
}
