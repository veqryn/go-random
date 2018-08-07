package random_test

import (
	"encoding/binary"
	"math"
	math_rand "math/rand"
	"testing"

	"github.com/veqryn/go-random"
)

func BenchmarkSecureRandSourceFloat64(b *testing.B) {
	b.ReportAllocs()
	source := math_rand.New(random.SecureRandSource)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		source.Float64()
	}
}

func BenchmarkSecureRandomStringBytesHex(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		random.SecureRandomStringBytes(benchmarkLength, random.HexBytes)
	}
}

func BenchmarkSecureRandomStringBytesAlphabet(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		random.SecureRandomStringBytes(benchmarkLength, random.AlphabetBytes)
	}
}

func BenchmarkSecureRandomStringBytesBase64(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		random.SecureRandomStringBytes(benchmarkLength, random.Base64URLBytes)
	}
}

func BenchmarkSecureRandomStringRunesAlphabet(b *testing.B) {
	b.ReportAllocs()
	runes := []rune(string(random.Alphabet))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.SecureRandomStringRunes(benchmarkLength, runes)
	}
}

func BenchmarkSecureRandomStringRunesBase64(b *testing.B) {
	b.ReportAllocs()
	runes := []rune(string(random.Base64URL))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.SecureRandomStringRunes(benchmarkLength, runes)
	}
}

func BenchmarkSecureRandomBitsLittleEndian(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.SecureRandomBits(8*benchmarkLength, binary.LittleEndian)
	}
}

func BenchmarkSecureRandomBitsBigEndian(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.SecureRandomBits(8*benchmarkLength, binary.BigEndian)
	}
}

func BenchmarkSecureRandomBitBlocks(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.SecureRandomBitBlocks(8*benchmarkLength, 5, binary.LittleEndian)
	}
}

func BenchmarkSecureRandomHexComparedToString(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.SecureRandomHex(benchmarkLength)
	}
}

func BenchmarkSecureRandomHexComparedToBytes(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.SecureRandomHex(2 * benchmarkLength)
	}
}

func BenchmarkSecureRandomBytes(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.SecureRandomBytes(benchmarkLength)
	}
}

func BenchmarkSecureRandomNumber(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.SecureRandomNumber(0, math.MaxInt64)
	}
}
