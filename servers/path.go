package servers

import (
	"errors"
	"regexp"
	"strings"
	"sync"
)

func NewPathTree(rootPath string) (pathTree *PathTree) {
	pathTree = &PathTree{
		pathCache: make(map[string]*PathTreeNode),
	}

	root := newPathTreeNode(pathTree, nil)
	if rootPath == "" {
		rootPath = "/"
	}
	root.path = rootPath
	root.key = rootPath[1:]
	pathTree.root = root
	return
}

func newPathTreeNode(tree *PathTree, parent *PathTreeNode) (node *PathTreeNode) {
	node = &PathTreeNode{
		keyCache: make(map[string]*PathTreeNode),
		parent:   parent,
		tree:     tree,
	}
	return
}

type PathTree struct {
	root          *PathTreeNode
	pathCache     map[string]*PathTreeNode // 通过 path 缓存
	pathCacheLock sync.Mutex
}

type PathTreeNode struct {
	key          string
	path         string
	keyCache     map[string]*PathTreeNode // 通过key缓存
	keyCacheLock sync.Mutex

	parent *PathTreeNode

	pathParamRules   []*PathParamRule
	pathParamRuleLen int
	matchRegexp      *regexp.Regexp
	matchRule        string

	hasMatchAll bool
	children    []*PathTreeNode

	extends []*PathTreeNodeExtend
	tree    *PathTree
}

type PathTreeNodeExtend struct {
	order  int
	extend interface{}
}

type PathParamRule struct {
	value  string
	isName bool
}

type PathMatchResult struct {
	Params     []*PathParam  `json:"params"`
	MatchRules []string      `json:"matchRules"`
	Node       *PathTreeNode `json:"node"`
}

type PathParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

var (
	defaultParamRule = `(.*)`
)

// AddPath 添加路径映射
// 普通: /x、/x/x
// 带参数：/x/{name}/xx、/x/{name}-{age}/xx
// 参数正则：/x/{name:a-Z0-9}/xx
// {name} 默认参数转为正则：(.*) 匹配任意字符
// 匹配任意路径：/x/x/{:**}，其中 ** 将匹配任意路径包含`/`字符也能匹配
// 顺序是 order 从小到大 执行
func (this_ *PathTree) AddPath(path string, order int, extend interface{}) (err error) {
	err = this_.root.add(strings.Split(path, ""), order, extend)

	return
}

func (this_ *PathTree) Match(path string) (matchList []*PathMatchResult, err error) {
	//find := this_.getByPath(path)
	//if find != nil {
	//	matchList = append(matchList, &PathMatchResult{
	//		Node: find,
	//	})
	//	//fmt.Println("match from path cache")
	//	//return
	//}

	matchList_, err := this_.root.match(strings.Split(path, ""))
	if err != nil {
		return
	}
	matchList = append(matchList, matchList_...)

	return
}

func (this_ *PathTreeNode) add(strList []string, order int, extend interface{}) (err error) {
	var strLen = len(strList)
	if strLen == 0 {
		err = errors.New("path cannot be empty")
		return
	}
	if strList[0] != "/" {
		err = errors.New("path must start with '/'")
		return
	}
	var matchRule string
	var pathParamRules []*PathParamRule
	var lastParamName string
	var lastParamRule string
	var lastParamOther string
	var isRuleStart bool
	var isPathParamStart bool
	var str string
	var strIndex = 1
	var curChar string
	for ; strIndex < strLen; strIndex++ {
		curChar = strList[strIndex]
		if curChar == "/" {
			break
		}
		str += strList[strIndex]

		if curChar == "{" {
			isPathParamStart = true
			continue
		}

		if curChar == "}" {
			if !isPathParamStart && !isRuleStart {
				err = errors.New("path [" + this_.path + strings.Join(strList, "") + "] analysis error, not found `{`")
				return
			}
			isPathParamStart = false
			isRuleStart = false

			if lastParamOther != "" {
				matchRule += "(" + lastParamOther + ")"
				pathParamRules = append(pathParamRules, &PathParamRule{
					value: lastParamOther,
				})
			}

			pathParamRules = append(pathParamRules, &PathParamRule{
				value:  lastParamName,
				isName: true,
			})
			lastParamName = ""
			if lastParamRule == "" {
				matchRule += defaultParamRule
			} else {
				if lastParamRule == "*" {
					matchRule += defaultParamRule
				} else {
					matchRule += "(" + lastParamRule + ")"
				}
			}
			lastParamOther = ""
			lastParamRule = ""
			continue
		}
		if isPathParamStart {
			if curChar == ":" {
				isPathParamStart = false
				isRuleStart = true
				continue
			}
			lastParamName += curChar
		} else if isRuleStart {
			lastParamRule += curChar
		} else {
			lastParamOther += curChar
		}
	}
	if isPathParamStart || isRuleStart {
		err = errors.New("path [" + this_.path + strings.Join(strList, "") + "] analysis error, not found `}`")
		return
	}
	if matchRule != "" {
		if lastParamOther != "" {
			matchRule += "(" + lastParamOther + ")"
			pathParamRules = append(pathParamRules, &PathParamRule{
				value: lastParamOther,
			})
		}
	}

	var hasMatchAll bool
	find := strings.Count(matchRule, "**")
	if find > 0 {
		if find > 1 {
			err = errors.New("path [" + this_.path + strings.Join(strList, "") + "] analysis error, has more `**`")
			return
		}
		if !strings.HasSuffix(matchRule, "(**)") {
			err = errors.New("path [" + this_.path + strings.Join(strList, "") + "] analysis error, must match end `(**)`")
			return
		}
		hasMatchAll = true
		matchRule = strings.ReplaceAll(matchRule, "**", ".*")
	}

	var key string

	if matchRule == "" {
		key = str
	} else {
		key = str
	}
	child := this_.getChild(key)

	var isNew = child == nil
	if child == nil {
		child = newPathTreeNode(this_.tree, this_)
		child.key = key
		if this_.parent == nil {
			child.path = "/" + str
		} else {
			child.path = this_.path + "/" + str
		}
		if matchRule != "" {
			matchRule = "^" + matchRule + "$"
			child.matchRule = matchRule
			child.matchRegexp, err = regexp.Compile(matchRule)
			if err != nil {
				err = errors.New("path [" + child.path + "] regexp compile error:" + err.Error())
				return
			}
		} else {
			child.matchRule = str
		}
		child.pathParamRules = pathParamRules
		child.pathParamRuleLen = len(pathParamRules)
		child.hasMatchAll = hasMatchAll
	}
	var nextStrList = strList[strIndex:]

	var isEnd = len(nextStrList) == 0
	if !isEnd {
		err = child.add(nextStrList, order, extend)
		if err != nil {
			return
		}
	} else {
		child.extends = append(child.extends, &PathTreeNodeExtend{
			order:  order,
			extend: extend,
		})
		//if !isNew {
		//	err = errors.New("path [" + child.Path + "] already exists")
		//	return
		//}
	}

	if isNew {
		this_.addChild(child)

		if isEnd {
			var hasMatch bool
			n := child
			for n != nil {
				if n.matchRegexp != nil {
					hasMatch = true
					break
				}
				n = n.parent
			}
			if !hasMatch {
				this_.tree.addPathCache(child.path, child)
			}
		}
	}

	return
}

func (this_ *PathTreeNode) GetExtends() []*PathTreeNodeExtend {
	return this_.extends
}

func (this_ *PathTreeNodeExtend) GetExtend() interface{} {
	return this_.extend
}

func (this_ *PathTreeNode) match(strList []string) (matchList []*PathMatchResult, err error) {
	var strLen = len(strList)
	if strLen == 0 {
		err = errors.New("path cannot be empty")
		return
	}
	if strList[0] != "/" {
		err = errors.New("path must start with '/'")
		return
	}
	var str string
	var curChar string
	var strIndex = 1
	for ; strIndex < strLen; strIndex++ {
		curChar = strList[strIndex]
		if curChar == "/" {
			break
		}
		str += strList[strIndex]
	}
	var nextStrList = strList[strIndex:]
	var matchList_ []*PathMatchResult
	for _, c := range this_.children {
		matchList_, err = childMatch(c, str, nextStrList)
		if err != nil {
			return
		}
		matchList = append(matchList, matchList_...)
	}
	return
}

func childMatch(child *PathTreeNode, matchStr string, nextStrList []string) (matchList []*PathMatchResult, err error) {
	var params []*PathParam
	if child.matchRegexp == nil {
		if child.matchRule != matchStr {
			return
		}
	} else {
		var matchRes [][]string
		if child.hasMatchAll {
			matchStr += strings.Join(nextStrList, "")
			nextStrList = []string{}
			matchRes = child.matchRegexp.FindAllStringSubmatch(matchStr, -1)
		} else {
			matchRes = child.matchRegexp.FindAllStringSubmatch(matchStr, -1)
		}

		//if child.HasMatchAll {
		//	fmt.Println("-----match start-----")
		//	fmt.Println("matchStr:", matchStr)
		//	fmt.Println("matchRule:", child.MatchRule)
		//	fmt.Println("pathParamRules:", util.GetStringValue(child.PathParamRules))
		//	fmt.Println("FindAllString:", util.GetStringValue(child.matchRegexp.FindAllString(matchStr, -1)))
		//	fmt.Println("FindAllStringSubMatch:", util.GetStringValue(child.matchRegexp.FindAllStringSubmatch(matchStr, -1)))
		//	fmt.Println("-----match end-----")
		//}
		var matched bool
		for _, finds := range matchRes {
			if len(finds) == 0 || finds[0] != matchStr {
				continue
			}
			vs := finds[1:]
			vsLen := len(vs)
			if vsLen != child.pathParamRuleLen {
				return
			}

			for i, v := range vs {
				pathParamRule := child.pathParamRules[i]

				if pathParamRule.isName {
					param := &PathParam{
						Name:  pathParamRule.value,
						Value: v,
					}
					params = append(params, param)
				} else if pathParamRule.value != v {
					return
				}
			}
			matched = true
			break
		}
		if !matched {
			return
		}
	}
	if len(nextStrList) > 0 {
		var matchList_ []*PathMatchResult
		matchList_, err = child.match(nextStrList)
		if err != nil {
			return
		}
		for _, match := range matchList_ {
			match.Params = append(params, match.Params...)
			match.MatchRules = append([]string{child.matchRule}, match.MatchRules...)
		}
		matchList = append(matchList, matchList_...)
	} else {
		matchList = append(matchList, &PathMatchResult{
			MatchRules: []string{child.matchRule},
			Node:       child,
			Params:     params,
		})
	}
	return
}

func (this_ *PathTreeNode) getChild(key string) (child *PathTreeNode) {
	this_.keyCacheLock.Lock()
	defer this_.keyCacheLock.Unlock()

	child = this_.keyCache[key]
	return
}

func (this_ *PathTreeNode) addChild(child *PathTreeNode) {
	this_.keyCacheLock.Lock()
	defer this_.keyCacheLock.Unlock()

	this_.keyCache[child.key] = child
	this_.children = append(this_.children, child)
	return
}

func (this_ *PathTree) getByPath(path string) (child *PathTreeNode) {
	this_.pathCacheLock.Lock()
	defer this_.pathCacheLock.Unlock()

	child = this_.pathCache[path]
	return
}

func (this_ *PathTree) addPathCache(path string, child *PathTreeNode) {
	this_.pathCacheLock.Lock()
	defer this_.pathCacheLock.Unlock()

	//fmt.Println("add path cache path:", path)
	this_.pathCache[path] = child
	return
}
