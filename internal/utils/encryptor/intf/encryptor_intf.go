package intf

type Encryptor interface {
	Encrypt() (encryptedBytes []byte, err error)

	Decrypt() (plainBytes []byte, err error)
}

type Signature interface {
	Hash() (fp string, err error)
}

type CoreSignContext struct {
	PlainBytes []byte
}

type CoreEncryptedContext struct {
	CoreSignContext
	EncryptBytes []byte
}
