package main

import (
	"github.com/miniclip/gonsul/errorutil"
	"os"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			var recoveredError = r.(errorutil.GonsulError)
			os.Exit(recoveredError.Code)
		}
	}()

	start()
}

