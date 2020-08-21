package logs

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

func isExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}

func makeLogDir(logPath string) error {
	logDir := path.Dir(logPath)
	if !isExist(logDir) {
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			return fmt.Errorf("create <%s> error: %s", logDir, err)
		}
	}

	return nil
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
