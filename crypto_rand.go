package random

import (
	"crypto/rand"
	"math/big"
)

func SecureRandomBytes(length int) []byte {

	var randomBytes = make([]byte, length)

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
