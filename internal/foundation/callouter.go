package foundation

import "fmt"

func CallOuter(autoRecover bool, reportError chan error, fun func()) (exception error) {
	if fun == nil {
		return nil
	}

	if autoRecover {
		defer func() {
			if info := recover(); info != nil {
				if err, ok := info.(error); ok {
					exception = err

					if reportError != nil {
						reportError <- exception
					}

				} else {
					exception = fmt.Errorf("%v", info)

					if reportError != nil {
						reportError <- exception
					}
				}
			}
		}()
	}

	fun()

	return
}
