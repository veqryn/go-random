package random_test

import (
	math_rand "math/rand"
	"strings"
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

func TestSecureRandomString(t *testing.T) {
	t.Parallel()
	for length := 0; length <= 128; length++ {
		result := random.SecureRandomString(length)
		if len(result) != length {
			t.Errorf("Expecting length %d; Got: %d", length, len(result))
		}
	}
}

func TestSecureRandomStringBytes(t *testing.T) {
	t.Parallel()
	for chars := 1; chars <= 256; chars++ {
		bytes := []byte(strings.Repeat("x", chars))
		for length := 0; length <= 128; length++ {
			result := random.SecureRandomStringBytes(length, bytes)
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
			result := random.SecureRandomStringRunes(length, runes)
			if len(result) != length {
				t.Errorf("Expecting length %d; Got: %d", length, len(result))
			}
		}
	}
}

func TestSecureRandSource(t *testing.T) {
	random.SecureRandSource.Uint64()
	random.SecureRandSource.Int63()
	random.SecureRandSource.Seed(1)
	math_rand.New(random.SecureRandSource)
}
