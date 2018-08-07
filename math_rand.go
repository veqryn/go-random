package random

import (
	"encoding/hex"
	"math"
	"math/bits"
	"math/rand"
	"strings"
)

// PseudoRandomString uses math/rand to return a random url-safe base64 string of given length.
// Uses the global math/rand instance, which locks on each call.
// Not cryptographically secure.
func PseudoRandomString(length int) string {
	return pseudoRandomStringBytesBase(rand.Uint64, length, Base64URLBytes)
}

// PseudoRandomStringRand uses math/rand to return a random url-safe base64 string of given length.
// Allows passing in rand source to avoid locking or to use other RNG's.
// Not cryptographically secure.
func PseudoRandomStringRand(rand *rand.Rand, length int) string {
	return pseudoRandomStringBytesBase(rand.Uint64, length, Base64URLBytes)
}

// PseudoRandomStringBytes uses math/rand to return a random string of given
// length made from the available character bytes.
// If the available character runes slice is empty, or length is negative, this will panic.
// This function is particularly efficient when the length of the availableCharRunes
// slice is a power of two.
// Uses the global math/rand instance, which locks on each call.
// Not cryptographically secure.
func PseudoRandomStringBytes(length int, availableCharBytes []byte) string {
	return pseudoRandomStringBytesBase(rand.Uint64, length, availableCharBytes)
}

// PseudoRandomStringBytesRand uses math/rand to return a random string of given
// length made from the available character bytes.
// If the available character runes slice is empty, or length is negative, this will panic.
// This function is particularly efficient when the length of the availableCharRunes
// slice is a power of two.
// Allows passing in rand source to avoid locking or to use other RNG's.
// Not cryptographically secure.
func PseudoRandomStringBytesRand(rand *rand.Rand, length int, availableCharBytes []byte) string {
	return pseudoRandomStringBytesBase(rand.Uint64, length, availableCharBytes)
}

// pseudoRandomStringBytesBase uses math/rand to return a random string of given
// length made from the available character bytes.
// If the available character runes slice is empty, or length is negative, this will panic.
// This function is particularly efficient when the length of the availableCharRunes
// slice is a power of two.
// Not cryptographically secure.
func pseudoRandomStringBytesBase(randUint64 func() uint64, length int, availableCharBytes []byte) string {

	// Check length
	if length < 0 {
		panic("random: length can not be negative")
	}

	availableCharLength := len(availableCharBytes)
	if availableCharLength == 0 {
		panic("random: availableCharBytes must not be empty")
	}

	// bitsNeeded is how many bits are needed to represent all available character options.
	// bitsNeeded is 1 less than the length because slices are zero based and the
	// highest bit value (which is the bitMask) would access the last index in the
	// available character slice (or beyond it slightly, which will be skipped).
	bitsNeeded := uint64(bits.Len64(uint64(availableCharLength) - 1))

	// If there is only 1 option
	if bitsNeeded == 0 || length == 0 {
		return strings.Repeat(string(availableCharBytes[0]), length)
	}

	// indicesPerUint64 is how many different letter indices can be found using a single uint64
	indicesPerUint64 := 64 / int(bitsNeeded)

	// bitMask is a mask (ie: 11111) that will allow for bitsNeededMaxLength permutations,
	// and will be used in bitwise operations against a random input to find the character index to use.
	var bitMask uint64 = 1<<bitsNeeded - 1

	// The resulting string
	result := make([]byte, length)
	completed := 0

	// Create the random string
	// Make call to retrieve random data on each iteration
	for randomBits := randUint64(); ; randomBits = randUint64() {

		// Cycle through blocks of random bits
		for attempted := 0; attempted < indicesPerUint64; attempted++ {

			// Mask bytes to get an index into the character slice
			charIdx := int(randomBits & bitMask)

			// Right shift to get rid of bits used
			randomBits >>= bitsNeeded

			// If charIdx is within availableCharLength, add that character to the random result string.
			// If not, we must ignore this randIdx in order to maintain equal probability and distribution.
			if charIdx < availableCharLength {
				result[completed] = availableCharBytes[charIdx]
				completed++
				if completed == length {
					return string(result)
				}
			}
		}
	}
}

// PseudoRandomStringRunes uses math/rand to return a random string of given
// length made from the available character runes.
// If the available character runes slice is empty, or length is negative, this will panic.
// This function is particularly efficient when the length of the availableCharRunes
// slice is a power of two.
// Uses the global math/rand instance, which locks on each call.
// Not cryptographically secure.
func PseudoRandomStringRunes(length int, availableCharRunes []rune) string {
	return pseudoRandomStringRunesBase(rand.Uint64, length, availableCharRunes)
}

// PseudoRandomStringRunesRand uses math/rand to return a random string of given
// length made from the available character runes.
// If the available character runes slice is empty, or length is negative, this will panic.
// This function is particularly efficient when the length of the availableCharRunes
// slice is a power of two.
// Allows passing in rand source to avoid locking or to use other RNG's.
// Not cryptographically secure.
func PseudoRandomStringRunesRand(rand *rand.Rand, length int, availableCharRunes []rune) string {
	return pseudoRandomStringRunesBase(rand.Uint64, length, availableCharRunes)
}

// pseudoRandomStringRunesBase uses math/rand to return a random string of given
// length made from the available character runes.
// If the available character runes slice is empty, or length is negative, this will panic.
// This function is particularly efficient when the length of the availableCharRunes
// slice is a power of two.
// Not cryptographically secure.
func pseudoRandomStringRunesBase(randUint64 func() uint64, length int, availableCharRunes []rune) string {

	// Check length
	if length < 0 {
		panic("random: length can not be negative")
	}

	availableCharLength := len(availableCharRunes)
	if availableCharLength == 0 {
		panic("random: availableCharRunes must not be empty")
	}

	// bitsNeeded is how many bits are needed to represent all available character options.
	// bitsNeeded is 1 less than the length because slices are zero based and the
	// highest bit value (which is the bitMask) would access the last index in the
	// available character slice (or beyond it slightly, which will be skipped).
	bitsNeeded := uint64(bits.Len64(uint64(availableCharLength) - 1))

	// If there is only 1 option
	if bitsNeeded == 0 || length == 0 {
		return strings.Repeat(string(availableCharRunes[0]), length)
	}

	// indicesPerUint64 is how many different letter indices can be found using a single uint64
	indicesPerUint64 := 64 / int(bitsNeeded)

	// bitMask is a mask (ie: 11111) that will allow for bitsNeededMaxLength permutations,
	// and will be used in bitwise operations against a random input to find the character index to use.
	var bitMask uint64 = 1<<bitsNeeded - 1

	// The resulting string
	result := make([]rune, length)
	completed := 0

	// Create the random string
	// Make call to retrieve random data on each iteration
	for randomBits := randUint64(); ; randomBits = randUint64() {

		// Cycle through blocks of random bits
		for attempted := 0; attempted < indicesPerUint64; attempted++ {

			// Mask bytes to get an index into the character slice
			charIdx := int(randomBits & bitMask)

			// Right shift to get rid of bits used
			randomBits >>= bitsNeeded

			// If charIdx is within availableCharLength, add that character to the random result string.
			// If not, we must ignore this randIdx in order to maintain equal probability and distribution.
			if charIdx < availableCharLength {
				result[completed] = availableCharRunes[charIdx]
				completed++
				if completed == length {
					return string(result)
				}
			}
		}
	}
}

// PseudoRandomBits uses math/rand to return the requested number of bits as a uint64 slice.
// Uses the global math/rand instance, which locks on each call.
// Not cryptographically secure.
func PseudoRandomBits(length int) []uint64 {

	// How long should the uint64 slice be?
	uint64Length := int(math.Ceil(float64(length) / 64))

	randomBits := make([]uint64, uint64Length)

	for i := 0; i < uint64Length; i++ {
		randomBits[i] = rand.Uint64()
	}
	return randomBits
}

// PseudoRandomBitsRand uses math/rand to return the requested number of bits as a uint64 slice.
// Allows passing in rand source to avoid locking or to use other RNG's.
// Not cryptographically secure.
func PseudoRandomBitsRand(rand *rand.Rand, length int) []uint64 {

	// How long should the uint64 slice be?
	uint64Length := int(math.Ceil(float64(length) / 64))

	randomBits := make([]uint64, uint64Length)

	for i := 0; i < uint64Length; i++ {
		randomBits[i] = rand.Uint64()
	}
	return randomBits
}

// PseudoRandomHex uses math/rand to return a slice of random hex data of a given length.
// Uses the global math/rand instance, which locks on each call.
// Not cryptographically secure.
func PseudoRandomHex(length int) string {
	return hex.EncodeToString(PseudoRandomBytes(length))
}

// PseudoRandomHexRand uses math/rand to return a slice of random hex data of a given length.
// Allows passing in rand source to avoid locking or to use other RNG's.
// Not cryptographically secure.
func PseudoRandomHexRand(rand *rand.Rand, length int) string {
	return hex.EncodeToString(PseudoRandomBytesRand(rand, length))
}

// PseudoRandomBytes uses math/rand to return the requested number of bytes.
// Uses the global math/rand instance, which locks on each call.
// Not cryptographically secure.
func PseudoRandomBytes(length int) []byte {
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		panic(err) // Impossible
	}
	return randomBytes
}

// PseudoRandomBytesRand uses math/rand to return the requested number of bytes.
// Allows passing in rand source to avoid locking or to use other RNG's.
// Not cryptographically secure.
func PseudoRandomBytesRand(rand *rand.Rand, length int) []byte {
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		panic(err) // Impossible
	}
	return randomBytes
}

// PseudoRandomInt63 uses math/rand to return a 63-bit number between [minInclusive, maxExclusive).
// Uses the global math/rand instance, which locks on each call.
// Not cryptographically secure. If max - min <= 0, this panics.
func PseudoRandomInt63(minInclusive, maxExclusive int64) int64 {
	return rand.Int63n(maxExclusive-minInclusive) + minInclusive
}

// PseudoRandomInt63Rand uses math/rand to return a 63-bit number between [minInclusive, maxExclusive).
// Allows passing in rand source to avoid locking or to use other RNG's.
// Not cryptographically secure. If max - min <= 0, this panics.
func PseudoRandomInt63Rand(rand *rand.Rand, minInclusive, maxExclusive int64) int64 {
	return rand.Int63n(maxExclusive-minInclusive) + minInclusive
}
