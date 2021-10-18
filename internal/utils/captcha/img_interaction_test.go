package captcha

import (
	"github.com/LittleBenx86/Benlog/internal/utils/collection"
	"github.com/LittleBenx86/Benlog/internal/utils/encryptor"
	"strings"
	"testing"
)

func Test_NewImageInteraction(t *testing.T) {
	originStr := "0123456789AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	RandomFontLib = strings.Split(originStr, "")
	err := collection.SliceShuffle(RandomFontLib)
	if err != nil {
		t.Fatal("random font lib create failed")
	}
	t.Log(RandomFontLib)

	it, err := NewImageInteraction("123456", 3, 6,
		collection.GenerateUnRepeatableRandomNumbers,
		func() func(plain string) (string, error) {
			return func(plain string) (string, error) {
				return encryptor.NewSHA256().SetPlainBytes([]byte(plain)).Hash()
			}
		}(),
	)
	if err != nil {
		t.Fatal("image interaction create failed")
	}
	t.Logf("%+v\n", it)
}
