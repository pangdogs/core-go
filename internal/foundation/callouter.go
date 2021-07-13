package foundation

import (
	"fmt"
	"runtime"
)

func CallOuter(autoRecover bool, reportError chan error, fun func()) (exception error) {
	if fun == nil {
		return nil
	}

	if autoRecover {
		defer func() {
			if info := recover(); info != nil {
				if reportError != nil {
					go func() {
						reportError <- ErrorAddStackTrace(info)
					}()
				}
			}
		}()
	}

	fun()

	return
}

func ErrorAddStackTrace(info interface{}) error {
	stackBuf := make([]byte, 4096)
	n := runtime.Stack(stackBuf, false)
	return fmt.Errorf("Error: %v\nStack: %s\n", info, stackBuf[:n])
}
