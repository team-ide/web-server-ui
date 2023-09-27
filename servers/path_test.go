package servers

import (
	"encoding/json"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"testing"
)

func TestPath(t *testing.T) {
	var err error

	pathTree := NewPathTree("/server")

	addPath := func(path string) {
		if err = pathTree.AddPath(path, nil); err != nil {
			panic(err)
		}
		return
	}

	addPath("/x")
	addPath("/x/x")
	addPath("/x/{age:[0-9]+}/xx")
	addPath("/x/{name}-{age}/xx")
	addPath("/x/{name}/xx")
	addPath("/b/a{:**}")

	bs, err := json.MarshalIndent(pathTree.Root, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))

	// 模拟 匹配

	var matchRes *PathMatchResult

	matchPath := func(path string) {
		if matchRes, err = pathTree.Match(path); err != nil {
			panic(err)
		}
		fmt.Println("match path:", path)
		fmt.Println("match result:", util.GetStringValue(matchRes))
	}

	matchPath("/x")
	matchPath("/x/x")
	matchPath("/x/")
	matchPath("/x/x/")
	matchPath("/x/张三")
	matchPath("/x/123")
	matchPath("/x/张三-123")
	matchPath("/x/张三-123/xxx")
	matchPath("/b/a张三-123/xxx")

}
