package handler

type ReviewService interface {
}

type Server struct {
	reviewService ReviewService
}

func NewServer() *Server {
	return &Server{}
}
