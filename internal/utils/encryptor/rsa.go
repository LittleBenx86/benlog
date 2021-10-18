package encryptor

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/LittleBenx86/Benlog/internal/utils/encryptor/intf"
	"log"
	"os"
	"runtime"
)

func GenerateRsaKey() error {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	x509PrivKey := x509.MarshalPKCS1PrivateKey(privKey)
	privFile, err := os.Create("benlog-priv.pem")
	if err != nil {
		return err
	}
	defer func(privFile *os.File) {
		err := privFile.Close()
		if err != nil {

		}
	}(privFile)

	privBlock := pem.Block{
		Type:  "BENLOG RSA PRIVATE KEY ",
		Bytes: x509PrivKey,
	}

	if err = pem.Encode(privFile, &privBlock); err != nil {
		return err
	}

	publicKey := privKey.PublicKey
	x509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		log.Fatal("unable to generate the public key from private key")
	}

	publicFile, _ := os.Create("benlog-public.pem")
	defer func(publicFile *os.File) {
		err := publicFile.Close()
		if err != nil {

		}
	}(publicFile)

	publicBlock := pem.Block{
		Type:  "BENLOG RSA PUBLIC KEY ",
		Bytes: x509PublicKey,
	}
	if err = pem.Encode(publicFile, &publicBlock); err != nil {
		return err
	}
	return nil
}

type RSA struct {
	intf.CoreEncryptedContext
	pubKey  []byte
	privKey []byte
}

func NewRSA() *RSA {
	return &RSA{}
}

func (r *RSA) SetPlainBytes(pb []byte) *RSA {
	r.PlainBytes = pb
	return r
}

func (r *RSA) SetEncryptedBytes(eb []byte) *RSA {
	r.EncryptBytes = eb
	return r
}

func (r *RSA) Encrypt() (encryptedBytes []byte, err error) {
	block, _ := pem.Decode(r.pubKey)
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				log.Printf("runtime error [%s], check the public key is valid\n", err)
			default:
				log.Printf("error [%s]\n", err)
			}
		}
	}()

	pkIntf, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return []byte{}, err
	}

	pk := pkIntf.(*rsa.PublicKey)
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, pk, r.PlainBytes)
	if err != nil {
		return []byte{}, err
	}
	return cipherText, nil
}

func (r *RSA) Decrypt() (plainBytes []byte, err error) {
	block, _ := pem.Decode(r.privKey)
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				log.Printf("runtime error [%s], check the public key is valid\n", err)
			default:
				log.Printf("error [%s]\n", err)
			}
		}
	}()

	pk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return []byte{}, err
	}

	plainBytes, err = rsa.DecryptPKCS1v15(rand.Reader, pk, r.EncryptBytes)
	if err != nil {
		return []byte{}, err
	}
	return plainBytes, nil
}
