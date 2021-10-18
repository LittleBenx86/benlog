package files_test

import (
	"os"
	"strings"
	"testing"

	"github.com/LittleBenx86/Benlog/internal/utils/files"
)

func Test_GetCurrentAbsolutePath(t *testing.T) {
	t.Log(files.GetCurrentFileExecAbsPath())
}

func Test_Getwd(t *testing.T) {
	var basePath string
	curPath, err := os.Getwd()
	t.Log(curPath)
	if err != nil {
		t.Fatal("get current work directory failed")
	}

	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test") {
		basePath = strings.Replace(strings.Replace(curPath, `\test`, "", 1), `/test`, "", 1)
	} else {
		basePath = curPath
	}
	t.Log(basePath)
}
