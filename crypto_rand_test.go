package random

import (
	math_rand "math/rand"
	"strings"
	"testing"
)

func BenchmarkSecureRandomStringBytes(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		SecureRandomStringBytes(40, AlphabetBytes)
	}
}

func BenchmarkSecureRandomStringRunes(b *testing.B) {
	b.ReportAllocs()
	runes := []rune(string(Alphabet))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SecureRandomStringRunes(40, runes)
	}
}

func TestSecureRandomStringBytes(t *testing.T) {
	t.Parallel()
	for chars := 1; chars <= 256; chars++ {
		bytes := []byte(strings.Repeat("x", chars))
		for length := 0; length <= 128; length++ {
			result := SecureRandomStringBytes(length, bytes)
			if len(result) != length {
				t.Errorf("Expecting length %d; Got: %d", length, len(result))
			}
		}
	}
}

func TestSecureRandomStringRunes(t *testing.T) {
	t.Parallel()
	for chars := 1; chars <= 300; chars++ {
		runes := []rune(strings.Repeat("x", chars))
		for length := 0; length <= 128; length++ {
			result := SecureRandomStringRunes(length, runes)
			if len(result) != length {
				t.Errorf("Expecting length %d; Got: %d", length, len(result))
			}
		}
	}
}

func TestSecureRandSource(t *testing.T) {
	SecureRandSource.Uint64()
	SecureRandSource.Int63()
	SecureRandSource.Seed(1)
	math_rand.New(SecureRandSource)
}
