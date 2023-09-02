package static

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	BaseDir string
)

func init() {
	var err error
	BaseDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	BaseDir, err = filepath.Abs(BaseDir)
	if err != nil {
		panic(err)
	}
	BaseDir = filepath.ToSlash(BaseDir)
	if !strings.HasSuffix(BaseDir, "/") {
		BaseDir += "/"
	}
}

// go test -timeout 3600s -v -run TestStatics ./statics
func TestStatic(t *testing.T) {
	fmt.Println("Dir:" + BaseDir)
	var err error
	var dist string
	dist, err = filepath.Abs(BaseDir + "../dist")
	if err != nil {
		panic(err)
	}
	dist = filepath.ToSlash(dist)
	err = SetAsset(dist, BaseDir, "dist")
	if err != nil {
		panic(err)
	}
}
