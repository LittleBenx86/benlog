package files

import (
	"fmt"
	"github.com/LittleBenx86/Benlog/internal/global/consts"
	"mime/multipart"
	"net/http"
	"os"

	InternalLogger "github.com/LittleBenx86/Benlog/internal/utils/logger"
)

// MIME, i.e. the Content-Type in header

func GetFilesMimeByFileName(path string) (string, error) {
	f, err := os.Open(path)
	defer f.Close()

	if err != nil {
		msg := fmt.Sprintf("%s, details: %s", consts.ERRORS_FILE_OPEN_ERR, err.Error())
		InternalLogger.GetInstance().Error(msg)
		return "", err
	}

	buf := make([]byte, 32)
	if _, err := f.Read(buf); err != nil {
		msg := fmt.Sprintf("%s, previous 32 bit unreadable, details: %s", consts.ERRORS_FILE_READ_ERR, err.Error())
		InternalLogger.GetInstance().Error(msg)
		return "", err
	}
	return http.DetectContentType(buf), nil
}

func GetFilesMimeByFileptr(fp multipart.File) (string, error) {
	buf := make([]byte, 32)
	if _, err := fp.Read(buf); err != nil {
		msg := fmt.Sprintf("%s, previous 32 bit unreadable, details: %s", consts.ERRORS_FILE_READ_ERR, err.Error())
		InternalLogger.GetInstance().Error(msg)
		return "", err
	}

	return http.DetectContentType(buf), nil
}
