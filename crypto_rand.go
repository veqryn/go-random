package random

import (
	"crypto/rand"
	"encoding/binary"
	"math"
	"math/big"
	"math/bits"
	math_rand "math/rand"
	"strings"
)

const (
	Hex                   = "0123456789abcdef"
	Alphabet              = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AlphabetUpperAndLower = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	AlphaNumeric          = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

var (
	HexBytes                   = []byte(Hex)
	AlphabetBytes              = []byte(Alphabet)
	AlphabetUpperAndLowerBytes = []byte(AlphabetUpperAndLower)
	AlphaNumericBytes          = []byte(AlphaNumeric)

	// SecureRandSource uses crypto/rand, is thread-safe, and implements math/rand.Source64.
	// To use, call math_rand.New(random.SecureRandSource) to get a *math/rand.Rand.
	SecureRandSource math_rand.Source64 = secureRandSource{}
)

// secureRandSource is an empty struct that implements math/rand.Source64
type secureRandSource struct{}

// Uint64 allows implementation of math/rand.Source64
func (s secureRandSource) Uint64() uint64 {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		panic("random: " + err.Error()) // Impossible
	}
	return binary.LittleEndian.Uint64(b)
}

// Uint63 allows implementation of math/rand.Source
func (s secureRandSource) Int63() int64 {
	return int64(s.Uint64() & ((1 << 63) - 1))
}

// Seed allows implementation of math/rand.Source
func (s secureRandSource) Seed(seed int64) {
	// no-op
}

// SecureRandomStringBytes uses crypto/rand to return a random string of given
// length made from the available character bytes.
// If the available character bytes slice is empty or greater than 256 in length, this will panic.
// This function is particularly efficient when the length of the availableCharBytes
// slice is a power of two.
func SecureRandomStringBytes(length int, availableCharBytes []byte) string {

	// Check lengths
	if length < 0 {
		panic("random: length can not be negative")
	}

	availableCharLength := len(availableCharBytes)
	if availableCharLength == 0 || availableCharLength > 256 {
		panic("random: availableCharBytes must not be empty or be longer than 256 bytes")
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

	// bitsNeededMaxLength is how many options could be represented max by bitsNeeded.
	// It will always be greater than or equal to the length of the available character options.
	var bitsNeededMaxLength uint64 = 1 << bitsNeeded

	// If the number of bits needed per index divides evenly into 8 (a byte),
	// and the available characters is equal to the bitMask's permutations (no overflow),
	// then the rest is quite simple and a lot of logic can be short circuited.
	if 8%bitsNeeded == 0 && int(bitsNeededMaxLength) == availableCharLength {
		return secureRandomStringBytesSimple(length, availableCharBytes, uint8(bitsNeeded), uint8(bitsNeededMaxLength))
	}

	// Otherwise things get fun
	return secureRandomStringBytesComplex(length, availableCharLength, availableCharBytes, bitsNeeded, bitsNeededMaxLength)
}

// secureRandomStringBytesSimple uses crypto/rand to return a random string of given
// length made from the available character bytes.
// This function can be used if the number of bits divides evenly into 8 (a byte)
// and the available characters is equal to the bitMask's permutations (no overflow).
func secureRandomStringBytesSimple(length int, availableCharBytes []byte, bitsNeeded uint8, bitsNeededMaxLength uint8) string {

	// indicesPerUint64 is how many different letter indices can be found using a single uint64
	indicesPerUint8 := 8 / int(bitsNeeded)

	// bitMask is a mask (ie: 11111) that will allow for bitsNeededMaxLength permutations,
	// and will be used in bitwise operations against a random input to find the character index to use.
	bitMask := bitsNeededMaxLength - 1

	// The resulting string
	result := make([]byte, length)

	// Make call to retrieve crypto/rand data
	randomBytes := SecureRandomBytes(int(math.Ceil(float64(length) / float64(indicesPerUint8))))

	// Create the random string
	for attempted := 0; attempted < length; attempted++ {

		// Find which index of random data to use
		randIdx := attempted / indicesPerUint8

		// Mask bits to get an index into the character slice
		charIdx := int(randomBytes[randIdx] & bitMask)

		// Right shift to get rid of bits used
		randomBytes[randIdx] >>= bitsNeeded

		// Put the byte at this index into the result
		result[attempted] = availableCharBytes[charIdx]
	}
	return string(result)
}

// secureRandomStringBytesSimple uses crypto/rand to return a random string of given
// length made from the available character bytes.
// This function uses uint64 slices of random data in order to decrease wasted bits, and will
// effectively deal with bit mask overflow while maintaining equal probability and distribution.
func secureRandomStringBytesComplex(length, availableCharLength int, availableCharBytes []byte, bitsNeeded, bitsNeededMaxLength uint64) string {

	// indicesPerUint64 is how many different letter indices can be found using a single uint64
	indicesPerUint64 := 64 / int(bitsNeeded)

	// bitMask is a mask (ie: 11111) that will allow for bitsNeededMaxLength permutations,
	// and will be used in bitwise operations against a random input to find the character index to use.
	bitMask := bitsNeededMaxLength - 1

	// The resulting string
	result := make([]byte, length)
	completed := 0

	// Create the random string
	for overflowMultiplier := maskOverflowMultiplier; ; overflowMultiplier += 1.0 {

		// bitBufferSize is the length still needed multiplied by bits needed per character.
		// When the bitMask can potentially overflow the available character options,
		// increase bitBufferSize by a percentage of that potential, to minimize system calls.
		// For example, if there are 5 available characters, then our mask allows for 8 max characters,
		// which gives a 37.5% chance of a missed hit. So for a length desired of 20 with 5 characters
		// (mask is 3 bits), instead of getting 60 bits of random data, we might get 105 bits of random data.
		bitBufferSize := int(float64(length-completed) * float64(bitsNeeded) *
			((overflowMultiplier + 1.0) - overflowMultiplier*float64(availableCharLength)/float64(bitsNeededMaxLength)))

		// Make call to retrieve crypto/rand data
		randomBits, bitBlockCount := SecureRandomBitBlocks(bitBufferSize, int(bitsNeeded), binary.LittleEndian)

		// Cycle through blocks of random bits
		for attempted := 0; attempted < bitBlockCount; attempted++ {

			// Find which index of random data to use
			randIdx := attempted / indicesPerUint64

			// Mask bits to get an index into the character slice
			charIdx := int(randomBits[randIdx] & bitMask)

			// Right shift to get rid of bits used
			randomBits[randIdx] >>= bitsNeeded

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

// SecureRandomStringRunes uses crypto/rand to return a random string of given
// length made from the available character runes.
// If the available character bytes slice is empty, this will panic.
// This function is particularly efficient when the length of the availableCharRunes
// slice is a power of two.
func SecureRandomStringRunes(length int, availableCharRunes []rune) string {

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

	// bitsNeededMaxLength is how many options could be represented max by bitsNeeded.
	// It will always be greater than or equal to the length of the available character options.
	var bitsNeededMaxLength uint64 = 1 << uint64(bitsNeeded)

	// indicesPerUint64 is how many different letter indices can be found using a single uint64
	indicesPerUint64 := 64 / int(bitsNeeded)

	// bitMask is a mask (ie: 11111) that will allow for bitsNeededMaxLength permutations,
	// and will be used in bitwise operations against a random input to find the character index to use.
	bitMask := bitsNeededMaxLength - 1

	// The resulting string
	result := make([]rune, length)
	completed := 0

	// Create the random string
	for overflowMultiplier := maskOverflowMultiplier; ; overflowMultiplier += 1.0 {

		// bitBufferSize is the length still needed multiplied by bits needed per character.
		// When the bitMask can potentially overflow the available character options,
		// increase bitBufferSize by a percentage of that potential, to minimize system calls.
		// For example, if there are 5 available characters, then our mask allows for 8 max characters,
		// which gives a 37.5% chance of a missed hit. So for a length desired of 20 with 5 characters
		// (mask is 3 bits), instead of getting 60 bits of random data, we might get 105 bits of random data.
		bitBufferSize := int(float64(length-completed) * float64(bitsNeeded) *
			((overflowMultiplier + 1.0) - overflowMultiplier*float64(availableCharLength)/float64(bitsNeededMaxLength)))

		// Make call to retrieve crypto/rand data
		randomBits, bitBlockCount := SecureRandomBitBlocks(bitBufferSize, int(bitsNeeded), binary.LittleEndian)

		// Cycle through blocks of random bits
		for attempted := 0; attempted < bitBlockCount; attempted++ {

			// Find which index of random data to use
			randIdx := attempted / indicesPerUint64

			// Mask bytes to get an index into the character slice
			charIdx := int(randomBits[randIdx] & bitMask)

			// Right shift to get rid of bits used
			randomBits[randIdx] >>= bitsNeeded

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

// stringMaskOverflowMultiplier controls how much extra random data to pull with crypto/rand
// when the bit mask allows for more possibilities than there are available character options.
// It can be any number equal to or greater than 0. Zero means do not pull any extract random bytes,
// while a value of 1 would mean that if there is are 25% more possibilities allowed by the mask
// than there are available character options, pull 25% more random data.
// A value of 2 would mean pull 50% more random data.
// The point of pulling more random data than possibly needed is to minimize the number of
// system calls that are made with crypto/rand, since those are costly.
// The current value was determined as a compromise between benchmarks
// and not pulling too much more random data than needed.
const maskOverflowMultiplier float64 = 1.5

// SecureRandomBits uses crypto/rand to return a slice of uint64 filled with the requested
// number of random bits.
// Unless bitLength is a multiple of 64, then the final uint64 will not be completely full
// of random bits.
// The binary.ByteOrder argument determines how the crypto/rand bytes get put into the []uint64.
// Note that for the final uint64 in the slice, LittleEndian fills from the low bits
// (right side) first, while BigEndian fills from the high bits (left side) first.
func SecureRandomBits(bitLength int, order binary.ByteOrder) []uint64 {
	bitz, _ := SecureRandomBitBlocks(bitLength, 1, order)
	return bitz
}

// SecureRandomBitBlocks uses crypto/rand to return a slice of uint64 filled with random bit data,
// using the byte order specified, as well as the number of usable bit blocks contained total.
// usableBlockSize is the number of bits that will be consumed at a time
// (for example, if using a mask to consume 3 bits at a time).
// bitLength should be a multiple of both usableBlockSize and 8 (crypto/rand gives random byte data),
// otherwise slightly more bits will be returned than requested.
// For example, if 62 bitLength is requested, and the usableBlockSize is 5, then this function determines
// that 12 blocks of 5 bits (60 bits) can be contained in each uint64. For the remaining 12 bits, it
// rounds up to 15 bits (as a multiple of 5 usableBlockSize) then determines that it will need to pull
// 2 bytes of random data (they only come in bytes) from crypto/rand to fulfill the remainder.
// It will then return a []uint64 of length 2 containing 80 bits of random data, and also returns
// the integer 75, which is the number of bits that are usable as determined by usableBlockSize.
// The binary.ByteOrder argument determines how the crypto/rand bytes get put into the []uint64.
// Note that for the final uint64 in the slice, LittleEndian fills from the low bits
// (right side) first, while BigEndian fills from the high bits (left side) first.
func SecureRandomBitBlocks(bitLength, usableBlockSize int, order binary.ByteOrder) ([]uint64, int) {

	// indicesPerUint64 is how many different usable blocks of bits a single uint64 can contain
	indicesPerUint64 := 64 / usableBlockSize

	// fullIndexByteCount is how many uint64's will be filled completely with random bits
	fullIndexByteCount := bitLength / (usableBlockSize * indicesPerUint64)

	// remainingBitsNeeded is how many bits are remaining, if any
	remainingBitsNeeded := bitLength % (indicesPerUint64 * usableBlockSize)

	// remainingBlockBitsNeeded is how many bits will be needed as a multiple of usableBlockSize
	remainingBlockBitsNeeded := int(math.Ceil(float64(remainingBitsNeeded)/float64(usableBlockSize))) * usableBlockSize

	// remainderByteCount is how many additional bytes will be needed to fulfill remainingBlockBitsNeeded
	remainderByteCount := int(math.Ceil(float64(remainingBlockBitsNeeded) / 8.0))

	// byteLength is how many bytes of random data are needed
	byteLength := remainderByteCount + (8 * fullIndexByteCount)

	// uint64Length is how long the uint64 slice should be
	uint64Length := fullIndexByteCount + int(math.Ceil(float64(remainderByteCount*8)/64.0))

	// Create slices
	randomBytes := make([]byte, uint64Length*8)
	randomBits := make([]uint64, uint64Length)

	// Read only the portion of random data that is needed, not the full slice
	if _, err := rand.Read(randomBytes[:byteLength]); err != nil {
		panic("random: " + err.Error()) // Impossible
	}

	// Set the bits using ByteOrder
	// If `64 % usableBlockSize >= 8` it is possible but not practical to use fewer than 8 bytes per uint64
	for i := range randomBits {
		randomBits[i] = order.Uint64(randomBytes[8*i:])
	}

	// Return randomBits and the number of usable bit blocks it contains
	return randomBits, indicesPerUint64*fullIndexByteCount + (8 * remainderByteCount / usableBlockSize)
}

// SecureRandomBytes uses crypto/rand to return a slice of random byte data of a given length
func SecureRandomBytes(length int) []byte {
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		panic("random: " + err.Error()) // Impossible
	}
	return randomBytes
}

// SecureRandomNumber uses crypto/rand to return a number between [minInclusive, maxExclusive)
func SecureRandomNumber(minInclusive int64, maxExclusive int64) int64 {
	min := big.NewInt(minInclusive)
	r, err := rand.Int(rand.Reader, big.NewInt(0).Sub(big.NewInt(maxExclusive), min))
	if err != nil {
		panic("random: " + err.Error()) // Impossible
	}
	return r.Add(r, min).Int64()
}
