package servers

func (this_ *Server) RegisterFilter(path string, order int, filter HttpFilter) (err error) {

	err = this_.filterPathTree.AddPath(path, order, filter)
	if err != nil {
		return
	}

	return
}

func (this_ *Server) RegisterMapper(path string, order int, mapper HttpMapper) (err error) {

	err = this_.mapperPathTree.AddPath(path, order, mapper)
	if err != nil {
		return
	}

	return
}

func (this_ *Server) RegisterInterceptor(path string, order int, interceptor HttpInterceptor) (err error) {

	err = this_.interceptorPathTree.AddPath(path, order, interceptor)
	if err != nil {
		return
	}
	return
}
