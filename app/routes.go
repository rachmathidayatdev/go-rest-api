package app

func initRoutes(s *Server) {
	s.BaseCRUDController("book", s.bookController)
}
