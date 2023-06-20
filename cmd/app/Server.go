package app

import (
	"github.com/RakhimovAns/Alif/pkg/service"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	Mux             *mux.Router
	CustomerService *service.CustomerService
}

func NewServer(mux *mux.Router, customerSvc *service.CustomerService) *Server {
	return &Server{Mux: mux, CustomerService: customerSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.Mux.ServeHTTP(writer, request)
}
