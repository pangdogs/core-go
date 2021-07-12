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
				exception = ErrorAddStackTrace(info)

				if reportError != nil {
					reportError <- exception
				}
			}
		}()
	}

	fun()

	return
}

func ErrorAddStackTrace(info interface{}) error {
	stackBuf := make([]byte, 1<<16)
	runtime.Stack(stackBuf, false)
	return fmt.Errorf("Error: %v\nStack:\n%v\n", info, stackBuf)
}
