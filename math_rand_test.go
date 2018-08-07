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

func TestPseudoRandomBits(t *testing.T) {
	t.Parallel()
	for length := 0; length <= 300; length++ {
		result := random.PseudoRandomBits(length)
		if len(result) != int(math.Ceil(float64(length)/64.0)) {
			t.Errorf("Expecting length %d; Got: %d", length, len(result))
		}
	}
}

func TestPseudoRandomBitsRand(t *testing.T) {
	t.Parallel()
	source := rand.New(rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	for length := 0; length <= 300; length++ {
		result := random.PseudoRandomBitsRand(source, length)
		if len(result) != int(math.Ceil(float64(length)/64.0)) {
			t.Errorf("Expecting length %d; Got: %d", length, len(result))
		}
	}
}

func TestPseudoRandomHex(t *testing.T) {
	t.Parallel()
	for length := 0; length <= 300; length++ {
		result := random.PseudoRandomHex(length)
		if len(result) != length {
			t.Errorf("Expecting length %d; Got: %d", length, len(result))
		}
	}
}

func TestPseudoRandomHexRand(t *testing.T) {
	t.Parallel()
	source := rand.New(rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	for length := 0; length <= 300; length++ {
		result := random.PseudoRandomHexRand(source, length)
		if len(result) != length {
			t.Errorf("Expecting length %d; Got: %d", length, len(result))
		}
	}
}

func TestPseudoRandomBytes(t *testing.T) {
	t.Parallel()
	for length := 0; length <= 300; length++ {
		result := random.PseudoRandomBytes(length)
		if len(result) != length {
			t.Errorf("Expecting length %d; Got: %d", length, len(result))
		}
	}
}

func TestPseudoRandomBytesRand(t *testing.T) {
	t.Parallel()
	source := rand.New(rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	for length := 0; length <= 300; length++ {
		result := random.PseudoRandomBytesRand(source, length)
		if len(result) != length {
			t.Errorf("Expecting length %d; Got: %d", length, len(result))
		}
	}
}

func TestPseudoRandomInt63(t *testing.T) {
	t.Parallel()
	increment := int64(math.MaxInt64 / 10)
	for min := int64(math.MinInt64); min < math.MaxInt64-increment; min += increment {
		until := min + math.MaxInt64
		if until <= min {
			until = math.MaxInt64
		}
		for max := min + 1; max <= until-increment; max += increment {
			num := random.PseudoRandomInt63(min, max)
			if num < min || num >= max {
				t.Errorf("Expected number in [%d, %d); Got: %d", min, max, num)
			}
		}
	}
}

func TestPseudoRandomInt63Rand(t *testing.T) {
	t.Parallel()
	source := rand.New(rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	increment := int64(math.MaxInt64 / 10)
	for min := int64(math.MinInt64); min < math.MaxInt64-increment; min += increment {
		until := min + math.MaxInt64
		if until <= min {
			until = math.MaxInt64
		}
		for max := min + 1; max <= until-increment; max += increment {
			num := random.PseudoRandomInt63Rand(source, min, max)
			if num < min || num >= max {
				t.Errorf("Expected number in [%d, %d); Got: %d", min, max, num)
			}
		}
	}
}
