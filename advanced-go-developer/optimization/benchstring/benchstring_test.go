package benchstring

import "testing"

func BenchmarkAlignRight(b *testing.B) {
	const (
		alignMe    string = "ultra-left"
		maxLen     int    = 15
		leftSymbol rune   = ' '
	)

	b.Run("naive", func(b *testing.B) {
		AlignRightNaive(alignMe, maxLen, leftSymbol)
	})

	b.Run("buf", func(b *testing.B) {
		AlignRightBuf(alignMe, maxLen, leftSymbol)
	})

	b.Run("repeat", func(b *testing.B) {
		AlignRightRepeat(alignMe, maxLen, leftSymbol)
	})

}
