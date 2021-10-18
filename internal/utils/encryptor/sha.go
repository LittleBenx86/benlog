package encryptor

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/LittleBenx86/Benlog/internal/utils/encryptor/intf"
)

// SHA
// Secure Hash Algorithm
// SHA256 (256 bits), others: SHA-1, SHA-224, SHA-384, SHA-512
// Digital Signature Standard (DSS) & Digital Signature Algorithm (DSA)

type SHA256 struct {
	intf.CoreSignContext
}

func NewSHA256() *SHA256 {
	return &SHA256{}
}

func (s *SHA256) SetPlainBytes(pb []byte) *SHA256 {
	s.PlainBytes = pb
	return s
}

func (s *SHA256) Hash() (hash string, err error) {
	digest := sha256.New()
	digest.Write(s.PlainBytes)
	hash = hex.EncodeToString(digest.Sum([]byte{})) // Generate the hash value
	return hash, nil
}
