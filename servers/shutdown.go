package servers

func (this_ *Server) Shutdown() {

	if this_.webListener != nil {
		_ = this_.webListener.Close()
	}
}
