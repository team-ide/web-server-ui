package servers

import "sort"

type PathMatchExtend struct {
	Params []*PathParam `json:"params"`
	Order  int          `json:"order"`
	Extend interface{}
}

func (this_ *Server) matchTree(path string, matchTree *PathTree, excludeTree *PathTree) (pathMatchExtends []*PathMatchExtend, err error) {

	matchList_, err := matchTree.Match(path)
	if err != nil {
		return
	}
	excludeList, err := excludeTree.Match(path)
	if err != nil {
		return
	}
	var excludeMap = make(map[*PathMatchResult]bool)
	for _, one := range excludeList {
		excludeMap[one] = true
	}
	var setMap = make(map[*PathMatchResult]bool)
	for _, one := range matchList_ {
		if excludeMap[one] {
			continue
		}
		if setMap[one] {
			continue
		}
		setMap[one] = true
		es := one.Node.GetExtends()
		for _, e := range es {
			pathMatchExtend := &PathMatchExtend{
				Params: one.Params,
				Order:  e.order,
				Extend: e.extend,
			}
			pathMatchExtends = append(pathMatchExtends, pathMatchExtend)
		}
	}

	// Order 正序
	sort.Slice(pathMatchExtends, func(i, j int) bool {
		return pathMatchExtends[i].Order < pathMatchExtends[j].Order
	})
	return
}
