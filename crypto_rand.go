package random

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"strconv"
)

func SecureRandomStringChars(length int, availableCharBytes []byte) string {

	// Check length
	availableCharLength := len(availableCharBytes)
	if availableCharLength == 0 || availableCharLength > math.MaxInt64 {
		panic("availableCharBytes length must be greater than 0 and less than or equal to 9223372036854775807")
	}

	// bitsNeeded is how many bits are needed to represent all available character options.
	// bitsNeeded is 1 less than the length because slices are zero based and the
	// highest bit value (which is the bitMask) would access the last index in the
	// available character slice (or beyond it slightly, which will be skipped).
	bitsNeeded := uint64(bits.Len64(uint64(availableCharLength) - 1))

	// bitsNeededMaxLength is how many options could be represented max by bitsNeeded.
	// It will always be greater than or equal to the length of the available character options.
	var bitsNeededMaxLength uint64 = 1 << uint64(bitsNeeded)

	// indicesPerUint64 is how many different letter indices can be found using a single uint64
	indicesPerUint64 := 64 / int(bitsNeeded)
	fmt.Println("init", length, availableCharLength, bitsNeeded, bitsNeededMaxLength, indicesPerUint64)

	// bitMask is a mask (ie: 11111) that will allow for bitsNeededOptionsLength permutations,
	// and will be used in bitwise operations against a random input.
	bitMask := bitsNeededMaxLength - 1
	fmt.Println("mask", strconv.FormatInt(int64(bitMask), 2))

	// The resulting string
	result := make([]byte, length)
	completed := 0

	// Create the random string
	for {

		// bitBufferSize is the length still needed times bits needed per character.
		// Increase bitBufferSize when there is overflow potential by double that potential, to minimize system calls.
		// For example, if there are 5 available characters, then our mask allows for 8 max characters, which gives 37.5% chance of a missed hit.
		// So for length desired of 20 with 5 characters, instead of getting 60 bits of random data, get 105 bits of random data.
		bitBufferSize := int(float64(length-completed) * float64(bitsNeeded) * (3.0 - 2.0*float64(availableCharLength)/float64(bitsNeededMaxLength)))
		bitBlockCount := int(math.Ceil(float64(bitBufferSize) / float64(bitsNeeded)))
		fmt.Println("buf", bitBufferSize)

		randomBits := SecureRandomBits(bitBufferSize, int(bitsNeeded), binary.LittleEndian)
		bufferSize := len(randomBits)
		for _, b := range randomBits {
			fmt.Printf("%064s ", strconv.FormatUint(uint64(b), 2))
		}
		fmt.Println("\nbuf size", bufferSize)

		// Find which uint64 in randomBits to use
		for attempted := 0; attempted < bitBlockCount; attempted++ {

			// Find random byte index
			randIdx := attempted / indicesPerUint64

			// Mask bytes to get an index into the character slice
			charIdx := int(randomBits[randIdx] & bitMask)
			fmt.Println("work", attempted, randIdx, strconv.FormatUint(uint64(randomBits[randIdx]), 2), strconv.FormatUint(uint64(charIdx), 2))

			// Right shift over the uint64 to get rid of bits used
			randomBits[randIdx] >>= bitsNeeded

			// If charIdx is within availableCharLength, add that character to the random string.
			// If not, we must ignore this randIdx in order to maintain equal probability.
			if charIdx < availableCharLength && charIdx%2 == 0 {
				result[completed] = availableCharBytes[charIdx]
				completed++
				fmt.Println("good", attempted, randIdx, strconv.FormatUint(uint64(randomBits[randIdx]), 2), strconv.FormatUint(uint64(charIdx), 2), string(availableCharBytes[charIdx]))
				if completed == length {
					return string(result)
				}
			}
		}

		//for randIdx := 0; randIdx < bufferSize ; randIdx++ {
		//
		//	for attempted := 0; attempted < indicesPerUint64; attempted++ {
		//
		//		// Mask bytes to get an index into the character slice
		//		charIdx := int(randomBits[randIdx] & bitMask)
		//		fmt.Println("work", randIdx, strconv.FormatUint(uint64(randomBits[randIdx]), 2), strconv.FormatUint(uint64(charIdx), 2))
		//
		//		// Right shift over the uint64 to get rid of bits used
		//		randomBits[randIdx] >>= bitsNeeded
		//
		//		// If charIdx is within availableCharLength, add that character to the random string.
		//		// If not, we must ignore this randIdx in order to maintain equal probability.
		//		if charIdx < availableCharLength && charIdx%2 == 0 {
		//			result[completed] = availableCharBytes[charIdx]
		//			completed++
		//			fmt.Println("good", randIdx, strconv.FormatUint(uint64(randomBits[randIdx]), 2), strconv.FormatUint(uint64(charIdx), 2), string(availableCharBytes[charIdx]))
		//			if completed == length {
		//				return string(result)
		//			}
		//		}
		//	}
		//}
	}
}

// usableBlockSize is number of bits that will be consumed at a time (for example, if using a mask to consume 3 bits a time).
// bitLength should be a multiple of usableBlockSize, otherwise more bits will be returned than requested.
func SecureRandomBits(bitLength, usableBlockSize int, order binary.ByteOrder) []uint64 {

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

	// TODO: if usableBitsPer64 is less than 56, it is
	// possible to use fewer than 8 bytes (of random data) per uint64
	for i := range randomBits {
		randomBits[i] = order.Uint64(randomBytes[8*i:])
	}
	return randomBits
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
