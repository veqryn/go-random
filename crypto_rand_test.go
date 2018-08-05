package random

import (
	"testing"
)

var (
	runes1  = []rune("0")
	runes2  = []rune("01")
	runes10 = []rune("0123456789")
	runes16 = []rune("0123456789abcdef")
	runes26 = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	runes32 = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ012345")
	runes52 = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	runes62 = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	runes64 = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")
	bytes1  = []byte(string(runes1))
	bytes2  = []byte(string(runes2))
	bytes10 = []byte(string(runes10))
	bytes16 = []byte(string(runes16))
	bytes26 = []byte(string(runes26))
	bytes32 = []byte(string(runes32))
	bytes52 = []byte(string(runes52))
	bytes62 = []byte(string(runes62))
	bytes64 = []byte(string(runes64))
)

func BenchmarkNaiveSecureRandomStringBytes(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		naiveSecureRandomStringChars(40, string(bytes26))
	}
}

func BenchmarkSecureRandomStringBytes2(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		SecureRandomStringBytes2(40, bytes26)
	}
}

func BenchmarkSecureRandomStringBytes(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		SecureRandomStringBytes(40, bytes26)
	}
}

func BenchmarkSecureRandomStringRunes(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		SecureRandomStringRunes(40, runes26)
	}
}

func TestSecureRandomStringBytes2(t *testing.T) {
	t.Parallel()
	t.Log(SecureRandomStringBytes2(64, bytes16))
}

func TestSecureRandomStringBytes(t *testing.T) {
	t.Parallel()
	t.Log(SecureRandomStringBytes(64, bytes1))
	t.Log(SecureRandomStringBytes(64, bytes2))
	t.Log(SecureRandomStringBytes(64, bytes10))
	t.Log(SecureRandomStringBytes(64, bytes16))
	t.Log(SecureRandomStringBytes(64, bytes26))
	t.Log(SecureRandomStringBytes(64, bytes32))
	t.Log(SecureRandomStringBytes(64, bytes52))
	t.Log(SecureRandomStringBytes(64, bytes62))
	t.Log(SecureRandomStringBytes(64, bytes64))
}

func TestSecureRandomStringRunes(t *testing.T) {
	t.Parallel()
	t.Log(SecureRandomStringRunes(64, runes1))
	t.Log(SecureRandomStringRunes(64, runes2))
	t.Log(SecureRandomStringRunes(64, runes10))
	t.Log(SecureRandomStringRunes(64, runes16))
	t.Log(SecureRandomStringRunes(64, runes26))
	t.Log(SecureRandomStringRunes(64, runes32))
	t.Log(SecureRandomStringRunes(64, runes52))
	t.Log(SecureRandomStringRunes(64, runes62))
	t.Log(SecureRandomStringRunes(64, runes64))
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
