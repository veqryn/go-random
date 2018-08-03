package random

import (
	"crypto/rand"
	"encoding/binary"
	"math"
	"math/big"
)

func SecureRandomBits(length int, order binary.ByteOrder) []uint64 {

	// How many bytes of random data are needed?
	byteLength := int(math.Ceil(float64(length) / 8))

	// How long should the uint64 slice be?
	uint64Length := int(math.Ceil(float64(length) / 64))

	randomBytes := make([]byte, uint64Length*8)
	randomBits := make([]uint64, uint64Length)

	// Read only the portion of random data that is needed, not the full slice
	_, err := rand.Read(randomBytes[:byteLength])
	if err != nil {
		panic(err) // Impossible
	}

	if order == binary.BigEndian {
		for idx := 0; idx < byteLength; idx += 8 {
			randomBits[idx/8] = binary.BigEndian.Uint64(randomBytes[idx : idx+8])
		}
		return randomBits
	}

	if order == binary.LittleEndian {
		for uint64Idx, byteIdx := 0, uint64Length*8; byteIdx > 0; uint64Idx, byteIdx = uint64Idx+1, byteIdx-8 {
			randomBits[uint64Idx] = binary.LittleEndian.Uint64(randomBytes[byteIdx-8 : byteIdx])
		}
		return randomBits
	}

	panic("Unrecognized binary.ByteOrder")
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
