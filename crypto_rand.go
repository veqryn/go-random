package random

import (
	"crypto/rand"
	"encoding/binary"
	"math"
	"math/big"
	"math/bits"
	"strings"
)

func SecureRandomStringBytes2(length int, availableCharBytes []byte) string {

	// Check length
	availableCharLength := len(availableCharBytes)
	if availableCharLength == 0 {
		panic("availableCharBytes must not be empty")
	}

	// bitsNeeded is how many bits are needed to represent all available character options.
	// bitsNeeded is 1 less than the length because slices are zero based and the
	// highest bit value (which is the bitMask) would access the last index in the
	// available character slice (or beyond it slightly, which will be skipped).
	bitsNeeded := uint64(bits.Len64(uint64(availableCharLength) - 1))

	// bitsNeededMaxLength is how many options could be represented max by bitsNeeded.
	// It will always be greater than or equal to the length of the available character options.
	var bitsNeededMaxLength uint8 = 1 << uint8(bitsNeeded)

	// indicesPerUint64 is how many different letter indices can be found using a single uint64
	indicesPerUint64 := 8 / int(bitsNeeded)

	// bitMask is a mask (ie: 11111) that will allow for bitsNeededOptionsLength permutations,
	// and will be used in bitwise operations against a random input.
	bitMask := bitsNeededMaxLength - 1

	// The resulting string
	result := make([]byte, length)
	completed := 0

	// Create the random string
	for {

		// bitBufferSize is the length still needed times bits needed per character.
		// When the bitMask can potentially overflow the available character options,
		// increase bitBufferSize by double or more of that potential, to minimize system calls.
		// For example, if there are 5 available characters, then our mask allows for 8 max characters,
		// which gives 37.5% chance of a missed hit. So for length desired of 20 with 5 characters
		// (mask is 3 bits), instead of getting 60 bits of random data, get 105 bits of random data.
		randomBits := SecureRandomBytes(length / indicesPerUint64)

		// Find which uint64 in randomBits to use
		for attempted := 0; attempted < length; attempted++ {

			// Find random byte index
			randIdx := attempted / indicesPerUint64

			// Mask bytes to get an index into the character slice
			charIdx := int(randomBits[randIdx] & bitMask)

			// Right shift over the uint64 to get rid of bits used
			randomBits[randIdx] >>= bitsNeeded

			// If charIdx is within availableCharLength, add that character to the random string.
			// If not, we must ignore this randIdx in order to maintain equal probability.
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

func SecureRandomStringBytes(length int, availableCharBytes []byte) string {

	// Check length
	availableCharLength := len(availableCharBytes)
	if availableCharLength == 0 {
		panic("availableCharBytes must not be empty")
	}

	// bitsNeeded is how many bits are needed to represent all available character options.
	// bitsNeeded is 1 less than the length because slices are zero based and the
	// highest bit value (which is the bitMask) would access the last index in the
	// available character slice (or beyond it slightly, which will be skipped).
	bitsNeeded := uint64(bits.Len64(uint64(availableCharLength) - 1))

	// If there is only 1 option
	if bitsNeeded == 0 {
		return strings.Repeat(string(availableCharBytes[0]), length)
	}

	// bitsNeededMaxLength is how many options could be represented max by bitsNeeded.
	// It will always be greater than or equal to the length of the available character options.
	var bitsNeededMaxLength uint64 = 1 << uint64(bitsNeeded)

	// indicesPerUint64 is how many different letter indices can be found using a single uint64
	indicesPerUint64 := 64 / int(bitsNeeded)

	// bitMask is a mask (ie: 11111) that will allow for bitsNeededOptionsLength permutations,
	// and will be used in bitwise operations against a random input.
	bitMask := bitsNeededMaxLength - 1

	// The resulting string
	result := make([]byte, length)
	completed := 0

	// Create the random string
	for overflowMultiplier := stringMaskOverflowMultiplier; ; overflowMultiplier += 1.0 {

		// bitBufferSize is the length still needed times bits needed per character.
		// When the bitMask can potentially overflow the available character options,
		// increase bitBufferSize by double or more of that potential, to minimize system calls.
		// For example, if there are 5 available characters, then our mask allows for 8 max characters,
		// which gives 37.5% chance of a missed hit. So for length desired of 20 with 5 characters
		// (mask is 3 bits), instead of getting 60 bits of random data, get 105 bits of random data.
		bitBufferSize := int(float64(length-completed) * float64(bitsNeeded) *
			((overflowMultiplier + 1.0) - overflowMultiplier*float64(availableCharLength)/float64(bitsNeededMaxLength)))

		randomBits, bitBlockCount := SecureRandomBitBlocks(bitBufferSize, int(bitsNeeded), binary.LittleEndian)

		// Find which uint64 in randomBits to use
		for attempted := 0; attempted < bitBlockCount; attempted++ {

			// Find random byte index
			randIdx := attempted / indicesPerUint64

			// Mask bytes to get an index into the character slice
			charIdx := int(randomBits[randIdx] & bitMask)

			// Right shift over the uint64 to get rid of bits used
			randomBits[randIdx] >>= bitsNeeded

			// If charIdx is within availableCharLength, add that character to the random string.
			// If not, we must ignore this randIdx in order to maintain equal probability.
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

func SecureRandomStringRunes(length int, availableCharRunes []rune) string {

	// Check length
	availableCharLength := len(availableCharRunes)
	if availableCharLength == 0 {
		panic("availableCharBytes must not be empty")
	}

	// bitsNeeded is how many bits are needed to represent all available character options.
	// bitsNeeded is 1 less than the length because slices are zero based and the
	// highest bit value (which is the bitMask) would access the last index in the
	// available character slice (or beyond it slightly, which will be skipped).
	bitsNeeded := uint64(bits.Len64(uint64(availableCharLength) - 1))

	// If there is only 1 option
	if bitsNeeded == 0 {
		return strings.Repeat(string(availableCharRunes[0]), length)
	}

	// bitsNeededMaxLength is how many options could be represented max by bitsNeeded.
	// It will always be greater than or equal to the length of the available character options.
	var bitsNeededMaxLength uint64 = 1 << uint64(bitsNeeded)

	// indicesPerUint64 is how many different letter indices can be found using a single uint64
	indicesPerUint64 := 64 / int(bitsNeeded)

	// bitMask is a mask (ie: 11111) that will allow for bitsNeededOptionsLength permutations,
	// and will be used in bitwise operations against a random input.
	bitMask := bitsNeededMaxLength - 1

	// The resulting string
	result := make([]rune, length)
	completed := 0

	// Create the random string
	for overflowMultiplier := stringMaskOverflowMultiplier; ; overflowMultiplier += 1.0 {

		// bitBufferSize is the length still needed times bits needed per character.
		// When the bitMask can potentially overflow the available character options,
		// increase bitBufferSize by double or more of that potential, to minimize system calls.
		// For example, if there are 5 available characters, then our mask allows for 8 max characters,
		// which gives 37.5% chance of a missed hit. So for length desired of 20 with 5 characters
		// (mask is 3 bits), instead of getting 60 bits of random data, get 105 bits of random data.
		bitBufferSize := int(float64(length-completed) * float64(bitsNeeded) *
			((overflowMultiplier + 1.0) - overflowMultiplier*float64(availableCharLength)/float64(bitsNeededMaxLength)))

		randomBits, bitBlockCount := SecureRandomBitBlocks(bitBufferSize, int(bitsNeeded), binary.LittleEndian)

		// Find which uint64 in randomBits to use
		for attempted := 0; attempted < bitBlockCount; attempted++ {

			// Find random byte index
			randIdx := attempted / indicesPerUint64

			// Mask bytes to get an index into the character slice
			charIdx := int(randomBits[randIdx] & bitMask)

			// Right shift over the uint64 to get rid of bits used
			randomBits[randIdx] >>= bitsNeeded

			// If charIdx is within availableCharLength, add that character to the random string.
			// If not, we must ignore this randIdx in order to maintain equal probability.
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
// It can be any number equal or greater than 0. Zero means do not pull any extract random bytes,
// while a value of 2 would mean that if there is are 25% more possibilities allowed by the mark
// than there are available character options, pull 50% more random data.
// The point of pulling more random data than possibly needed is to minimize the number of
// system calls that are made with crypto/rand, since those are costly.
var stringMaskOverflowMultiplier float64 = 1.5

func SecureRandomBits(bitLength int, order binary.ByteOrder) []uint64 {
	bitz, _ := SecureRandomBitBlocks(bitLength, 1, order)
	return bitz
}

// usableBlockSize is number of bits that will be consumed at a time (for example, if using a mask to consume 3 bits a time).
// bitLength should be a multiple of usableBlockSize, otherwise more bits will be returned than requested.
func SecureRandomBitBlocks(bitLength, usableBlockSize int, order binary.ByteOrder) ([]uint64, int) {

	// indicesPerUint64 is how many different usable blocks of bits a single uint64 can contain
	indicesPerUint64 := 64 / usableBlockSize

	// fullIndexByteCount is how many uint64's will be filled completely with random bits
	fullIndexByteCount := bitLength / (usableBlockSize * indicesPerUint64)

	// remainingBitsNeeded is how many bits are remaining if any
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
	_, err := rand.Read(randomBytes[:byteLength])
	if err != nil {
		panic(err) // Impossible
	}

	// TODO: if (64 % usableBlockSize) is less than 56, it is
	// possible to use fewer than 8 bytes (of random data) per uint64.
	for i := range randomBits {
		randomBits[i] = order.Uint64(randomBytes[8*i:])
	}

	// Return randomBits and the number of usable bit blocks it contains
	return randomBits, indicesPerUint64*fullIndexByteCount + (8 * remainderByteCount / usableBlockSize)
}

func SecureRandomBytes(length int) []byte {

	randomBytes := make([]byte, length)

	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err) // Impossible
	}

	return randomBytes
}

func SecureRandomNumber(minInclusive int64, maxExclusive int64) int64 {
	min := big.NewInt(minInclusive)
	r, err := rand.Int(rand.Reader, big.NewInt(0).Sub(big.NewInt(maxExclusive), min))
	if err != nil {
		panic(err)
	}
	return r.Add(r, min).Int64()
}
