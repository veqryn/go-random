// Package random ...
package random

const (
	Hex                   = "0123456789abcdef"
	Alphabet              = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AlphabetUpperAndLower = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	AlphaNumeric          = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	Base64URL             = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	Base64Std             = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var (
	HexBytes                   = []byte(Hex)
	AlphabetBytes              = []byte(Alphabet)
	AlphabetUpperAndLowerBytes = []byte(AlphabetUpperAndLower)
	AlphaNumericBytes          = []byte(AlphaNumeric)
	Base64URLBytes             = []byte(Base64URL)
	Base64StdBytes             = []byte(Base64Std)
)
