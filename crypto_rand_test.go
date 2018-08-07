package random_test

import (
	"encoding/binary"
	"math"
	math_rand "math/rand"
	"strings"
	"testing"

	"github.com/veqryn/go-random"
)

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

func TestSecureRandomBits(t *testing.T) {
	t.Parallel()
	for length := 0; length <= 300; length++ {
		result := random.SecureRandomBits(length, binary.LittleEndian)
		if len(result) != int(math.Ceil(float64(length)/64.0)) {
			t.Errorf("Expecting length %d; Got: %d", length, len(result))
		}
	}
	for length := 0; length <= 300; length++ {
		result := random.SecureRandomBits(length, binary.BigEndian)
		if len(result) != int(math.Ceil(float64(length)/64.0)) {
			t.Errorf("Expecting length %d; Got: %d", length, len(result))
		}
	}
}

func TestSecureRandomBitBlocks(t *testing.T) {
	t.Parallel()
	for length := 0; length <= 300; length++ {
		for blockSize := 1; blockSize <= 64; blockSize++ {
			result, blocks := random.SecureRandomBitBlocks(length, blockSize, binary.LittleEndian)
			if len(result) < int(math.Ceil(float64(length)/64.0)) {
				t.Errorf("Expecting length %d; Got: %d", int(math.Ceil(float64(length)/64.0)), len(result))
			}
			if blocks*blockSize < length {
				t.Errorf("Expecting blocks at least %d; Got: %d", length/64, blocks)
			}
		}
	}
	for length := 0; length <= 300; length++ {
		for blockSize := 1; blockSize <= 64; blockSize++ {
			result, blocks := random.SecureRandomBitBlocks(length, blockSize, binary.BigEndian)
			if len(result) < int(math.Ceil(float64(length)/64.0)) {
				t.Errorf("Expecting length %d; Got: %d", int(math.Ceil(float64(length)/64.0)), len(result))
			}
			if blocks*blockSize < length {
				t.Errorf("Expecting blocks at least %d; Got: %d", length, blocks)
			}
		}
	}
}

func TestSecureRandomHex(t *testing.T) {
	t.Parallel()
	for length := 0; length <= 300; length++ {
		result := random.SecureRandomHex(length)
		if len(result) != length {
			t.Errorf("Expecting length %d; Got: %d", length, len(result))
		}
	}
}

func TestSecureRandomBytes(t *testing.T) {
	t.Parallel()
	for length := 0; length <= 300; length++ {
		result := random.SecureRandomBytes(length)
		if len(result) != length {
			t.Errorf("Expecting length %d; Got: %d", length, len(result))
		}
	}
}

func TestSecureRandomNumber(t *testing.T) {
	t.Parallel()
	increment := int64(math.MaxInt64 / 200)
	source := math_rand.New(math_rand.NewSource(random.SecureRandomNumber(math.MinInt64, math.MaxInt64)))
	for min := int64(math.MinInt64); min < math.MaxInt64-increment; min += increment {
		max := random.PseudoRandomInt63Rand(source, 0, math.MaxInt64) - random.PseudoRandomInt63Rand(source, 0, math.MaxInt64)
		if max <= min {
			max = min + 1
		}
		num := random.SecureRandomNumber(min, max)
		if num < min || num >= max {
			t.Errorf("Expected number in [%d, %d); Got: %d", min, max, num)
		}
	}
}

func TestSecureRandSource(t *testing.T) {
	random.SecureRandSource.Uint64()
	random.SecureRandSource.Int63()
	random.SecureRandSource.Seed(1)
	math_rand.New(random.SecureRandSource).Float64()
}
