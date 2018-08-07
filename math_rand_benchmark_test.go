package random_test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/veqryn/go-random"
)

const benchmarkLength = 64

func BenchmarkPseudoRandomFloat64(b *testing.B) {
	b.ReportAllocs()
	source := rand.New(rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		source.Float64()
	}
}

func BenchmarkPseudoRandomStringBytesRandHex(b *testing.B) {
	b.ReportAllocs()
	source := rand.New(rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.PseudoRandomStringBytesRand(source, benchmarkLength, random.HexBytes)
	}
}

func BenchmarkPseudoRandomStringBytesRandAlphabet(b *testing.B) {
	b.ReportAllocs()
	source := rand.New(rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.PseudoRandomStringBytesRand(source, benchmarkLength, random.AlphabetBytes)
	}
}

func BenchmarkPseudoRandomStringBytesRandBase64(b *testing.B) {
	b.ReportAllocs()
	source := rand.New(rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.PseudoRandomStringBytesRand(source, benchmarkLength, random.Base64URLBytes)
	}
}

func BenchmarkPseudoRandomStringRunesRandAlphabet(b *testing.B) {
	b.ReportAllocs()
	r := []rune(random.Alphabet)
	source := rand.New(rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.PseudoRandomStringRunesRand(source, benchmarkLength, r)
	}
}

func BenchmarkPseudoRandomStringRunesRandBase64(b *testing.B) {
	b.ReportAllocs()
	r := []rune(random.Base64URL)
	source := rand.New(rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.PseudoRandomStringRunesRand(source, benchmarkLength, r)
	}
}

func BenchmarkPseudoRandomBitsRand(b *testing.B) {
	b.ReportAllocs()
	source := rand.New(rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.PseudoRandomBitsRand(source, 8*benchmarkLength)
	}
}

func BenchmarkPseudoRandomHexRandComparedToString(b *testing.B) {
	b.ReportAllocs()
	source := rand.New(rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.PseudoRandomHexRand(source, benchmarkLength)
	}
}

func BenchmarkPseudoRandomHexRandComparedToBytes(b *testing.B) {
	b.ReportAllocs()
	source := rand.New(rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.PseudoRandomHexRand(source, 2*benchmarkLength)
	}
}

func BenchmarkPseudoRandomBytesRand(b *testing.B) {
	b.ReportAllocs()
	source := rand.New(rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.PseudoRandomBytesRand(source, benchmarkLength)
	}
}

func BenchmarkPseudoRandomInt63Rand(b *testing.B) {
	b.ReportAllocs()
	source := rand.New(rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.PseudoRandomInt63Rand(source, 0, math.MaxInt64)
	}
}
