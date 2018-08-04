package random

import (
	"math"
	"math/rand"
)

// PseudoRandomBits returns the requested number of uint64 using math/rand (which is not cryptographically secure)
func PseudoRandomBits(rand *rand.Rand, length int) []uint64 {

	// How long should the uint64 slice be?
	uint64Length := int(math.Ceil(float64(length) / 64))

	randomBits := make([]uint64, uint64Length)

	for i := 0; i < uint64Length; i++ {
		randomBits[i] = rand.Uint64()
	}
	return randomBits
}

// PseudoRandomBytes returns the requested number of bytes using math/rand (which is not cryptographically secure)
func PseudoRandomBytes(rand *rand.Rand, length int) []byte {

	randomBytes := make([]byte, length)

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
