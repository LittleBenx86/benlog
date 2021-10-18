package encryptor_test

import (
	"encoding/base64"
	"fmt"
	"github.com/LittleBenx86/Benlog/internal/utils/encryptor"
	"testing"
)

func Test_GenerateRsaKey(t *testing.T) {
	err := encryptor.GenerateRsaKey()
	if err != nil {
		return
	}
}

func Test_Base64Decrypt(t *testing.T) {
	encryptedText := "666"
	decodeString, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		t.Logf("err [%s]\n", err)
	}
	fmt.Println(string(decodeString))
}

func Test_GenerateHash(t *testing.T) {
	h, _ := encryptor.NewSHA256().SetPlainBytes([]byte("666")).Hash()
	t.Log(h)
}
