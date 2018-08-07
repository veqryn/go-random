package random_test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/veqryn/go-random"
)

func BenchmarkPseudoRandomStringRunesRand(b *testing.B) {
	b.ReportAllocs()
	r := []rune(random.Alphabet)
	source := rand.New(rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.PseudoRandomStringRunesRand(source, 40, r)
	}
}

func BenchmarkPseudoRandomStringBytesRand(b *testing.B) {
	b.ReportAllocs()
	source := rand.New(rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.PseudoRandomStringBytesRand(source, 40, random.AlphabetBytes)
	}
}
