package commons

import (
	"fmt"
	"runtime"
)

var (
	version = "0.0.0"
)

func Version() string {
	return version
}

func OutVersion() {
	fmt.Println("Server Version : " + Version())
	fmt.Println("GO OS          : " + runtime.GOOS)
	fmt.Println("GO ARCH        : " + runtime.GOARCH)
	fmt.Println("GO Version     : " + runtime.Version())
}
