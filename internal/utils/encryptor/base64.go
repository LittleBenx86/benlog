package encryptor

import (
	"encoding/base64"
	"github.com/LittleBenx86/Benlog/internal/utils/encryptor/intf"
)

type Base64 struct {
	intf.CoreEncryptedContext
}

func NewBase64() *Base64 {
	return &Base64{}
}

func (b *Base64) SetPlainBytes(pb []byte) *Base64 {
	b.PlainBytes = pb
	return b
}

func (b *Base64) SetEncryptedBytes(eb []byte) *Base64 {
	b.EncryptBytes = eb
	return b
}

func (b *Base64) Encrypt() (encryptedBytes []byte, err error) {
	encryptedBytes = []byte(base64.StdEncoding.EncodeToString(b.PlainBytes))
	return encryptedBytes, nil
}

func (b *Base64) Decrypt() (plainBytes []byte, err error) {
	plainBytes, err = base64.StdEncoding.DecodeString(string(b.EncryptBytes))
	return
}
