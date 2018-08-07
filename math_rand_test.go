package random_test

import (
	"math"
	"math/rand"
	"strings"
	"testing"

	"github.com/veqryn/go-random"
)

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

func TestPseudoRandomHex(t *testing.T) {
	t.Parallel()
	random.PseudoRandomHex(0)
}

func TestPseudoRandomHexRand(t *testing.T) {
	t.Parallel()
	random.PseudoRandomHexRand(nil, 0)
}

func TestPseudoRandomBits(t *testing.T) {
	t.Parallel()
	random.PseudoRandomBits(0)
}

func TestPseudoRandomBitsRand(t *testing.T) {
	t.Parallel()
	random.PseudoRandomBitsRand(nil, 0)
}

func TestPseudoRandomBytes(t *testing.T) {
	t.Parallel()
	random.PseudoRandomBytes(0)
}

func TestPseudoRandomBytesRand(t *testing.T) {
	t.Parallel()
	random.PseudoRandomBytesRand(nil, 0)
}

func TestPseudoRandomInt63(t *testing.T) {
	t.Parallel()
	random.PseudoRandomInt63(0, 0)
}

func TestPseudoRandomInt63Rand(t *testing.T) {
	t.Parallel()
	random.PseudoRandomInt63Rand(nil, 0, 0)
}
