package files

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func GetProjectRuntimeRootPath() (rootPath string) {
	fileAbsPath := GetCurrentFileExecAbsPath()
	if strings.HasSuffix(fileAbsPath, `/internal/utils/files`) {
		rootPath = strings.Replace(fileAbsPath, `/internal/utils/files`, "", 1)
	} else {
		rootPath = fileAbsPath
	}
	return
}

// GetCurrentFileExecAbsPath
// Get the current executable application correct absolute path.
// Tips:
// Get absolute path by os.Executable() with go run and go build.
// go run: Run as local, it will compile the source codes to a temporary dir, which is the value of TEMP or TMP system environment variable,
//   then startup the application.
// go build: Only compile the project as an executable application binary file.
// go run main.go is equivalent to go build & ./main.
func GetCurrentFileExecAbsPath() (absPath string) {
	absPath = getCurrentAbsolutePathByExecutable()
	tmpPath, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(absPath, tmpPath) {
		absPath = getCurrentAbsolutePathByCaller()
	}
	return
}

// getCurrentAbsolutePathByExecutable
// Works for go build
func getCurrentAbsolutePathByExecutable() (absPath string) {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	absPath, _ = filepath.EvalSymlinks(filepath.Dir(execPath))
	return
}

// getCurrentAbsolutePathByCaller
// Works for go run
func getCurrentAbsolutePathByCaller() (absPath string) {
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		absPath = path.Dir(filename)
	}
	return
}
