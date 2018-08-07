package random_test

import (
	"testing"

	"github.com/veqryn/go-random"
)

func BenchmarkSecureRandomStringBytes(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		random.SecureRandomStringBytes(40, random.AlphabetBytes)
	}
}

func BenchmarkSecureRandomStringRunes(b *testing.B) {
	b.ReportAllocs()
	runes := []rune(string(random.Alphabet))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.SecureRandomStringRunes(40, runes)
	}
}
