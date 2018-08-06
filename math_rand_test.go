package random_test

import (
	"math"
	"math/rand"
	"strings"
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

func TestPseudoRandomString(t *testing.T) {
	t.Parallel()
	for length := 0; length <= 128; length++ {
		result := random.PseudoRandomString(length)
		if len(result) != length {
			t.Errorf("Expecting length %d; Got: %d", length, len(result))
		}
	}
}

func TestPseudoRandomStringRand(t *testing.T) {
	t.Parallel()
	source := rand.New(rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	for length := 0; length <= 128; length++ {
		result := random.PseudoRandomStringRand(source, length)
		if len(result) != length {
			t.Errorf("Expecting length %d; Got: %d", length, len(result))
		}
	}
}

func TestPseudoRandomStringBytes(t *testing.T) {
	t.Parallel()
	for chars := 1; chars <= 256; chars++ {
		bytes := []byte(strings.Repeat("x", chars))
		for length := 0; length <= 128; length++ {
			result := random.PseudoRandomStringBytes(length, bytes)
			if len(result) != length {
				t.Errorf("Expecting length %d; Got: %d", length, len(result))
			}
		}
	}
}

func TestPseudoRandomStringBytesRand(t *testing.T) {
	t.Parallel()
	source := rand.New(rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	for chars := 1; chars <= 256; chars++ {
		bytes := []byte(strings.Repeat("x", chars))
		for length := 0; length <= 128; length++ {
			result := random.PseudoRandomStringBytesRand(source, length, bytes)
			if len(result) != length {
				t.Errorf("Expecting length %d; Got: %d", length, len(result))
			}
		}
	}
}

func TestPseudoRandomStringRunes(t *testing.T) {
	t.Parallel()
	for chars := 1; chars <= 300; chars++ {
		runes := []rune(strings.Repeat("x", chars))
		for length := 0; length <= 128; length++ {
			result := random.PseudoRandomStringRunes(length, runes)
			if len(result) != length {
				t.Errorf("Expecting length %d; Got: %d", length, len(result))
			}
		}
	}
}

func TestPseudoRandomStringRunesRand(t *testing.T) {
	t.Parallel()
	source := rand.New(rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	for chars := 1; chars <= 300; chars++ {
		runes := []rune(strings.Repeat("x", chars))
		for length := 0; length <= 128; length++ {
			result := random.PseudoRandomStringRunesRand(source, length, runes)
			if len(result) != length {
				t.Errorf("Expecting length %d; Got: %d", length, len(result))
			}
		}
	}
}
