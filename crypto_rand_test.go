package random

import (
	"strings"
	"testing"
)

func BenchmarkNaiveSecureRandomStringBytes(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		naiveSecureRandomStringChars(40, string(Alphabet))
	}
}

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

func naiveSecureRandomStringChars(length int, availableCharBytes string) string {

	// Compute bitMask
	availableCharLength := len(availableCharBytes)
	if availableCharLength == 0 || availableCharLength > 256 {
		panic("availableCharBytes length must be greater than 0 and less than or equal to 256")
	}
	var bitLength byte
	var bitMask byte
	for bits := availableCharLength - 1; bits != 0; {
		bits = bits >> 1
		bitLength++
	}
	bitMask = 1<<bitLength - 1

	// Compute bufferSize: length + 2*(1 - availCharLength/bitMask+1)
	bufferSize := length + int(2.0*(1.0-(float64(availableCharLength)/float64(bitMask+1))))

	// Create random string
	result := make([]byte, length)
	for i, j, randomBytes := 0, 0, []byte{}; i < length; j++ {
		if j%bufferSize == 0 {
			// Random byte buffer is empty, get a new one
			randomBytes = SecureRandomBytes(bufferSize)
		}
		// Mask bytes to get an index into the character slice
		if idx := int(randomBytes[j%length] & bitMask); idx < availableCharLength {
			result[i] = availableCharBytes[idx]
			i++
		}
	}

	return string(result)
}
