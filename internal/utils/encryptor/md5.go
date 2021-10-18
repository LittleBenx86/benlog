package encryptor

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/LittleBenx86/Benlog/internal/utils/encryptor/intf"
)

// MD5
// Message-Digest Algorithm (hash)
// 1. irreversible
// 2. high degree of discreteness
// 3. compressibility (128 bits)
// 4. weak collision resistance
type MD5 struct {
	intf.CoreSignContext
}

func NewMD5() *MD5 {
	return &MD5{}
}

func (m *MD5) SetPlainBytes(pb []byte) *MD5 {
	m.PlainBytes = pb
	return m
}

func (m *MD5) Hash() (fp string, err error) {
	ctx := md5.New()                // md5 init
	ctx.Write(m.PlainBytes)         // md5 update
	cipher := ctx.Sum([]byte{})     // md5 final
	fp = hex.EncodeToString(cipher) // hex_digest
	return fp, nil
}
