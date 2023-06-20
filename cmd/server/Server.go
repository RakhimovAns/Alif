package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	Mux *mux.Router
}

func NewServer(mux *mux.Router) *Server {
	return &Server{Mux: mux}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.Mux.ServeHTTP(writer, request)
}
