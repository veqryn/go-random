package random

import "math/rand"

// PseudoRandomBytes returns the requested number of bytes using math/rand (which is not cryptographically secure)
func PseudoRandomBytes(rand *rand.Rand, length int) []byte {

	var randomBytes = make([]byte, length)

	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err) // Impossible
	}

	return randomBytes
}

// PseudoRandomNumber generates a random 63-bit number [min, max) between a inclusive minimum
// and an exclusive maximum, using math/rand (which is not cryptographically secure).
// If max - min <= 0, this panics.
func PseudoRandomNumber(rand *rand.Rand, minInclusive int64, maxExclusive int64) int64 {
	return rand.Int63n(maxExclusive-minInclusive) + maxExclusive
}
