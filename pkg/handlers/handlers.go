package handlers

import "github.com/RakhimovAns/Alif/cmd/app"

type Server struct {
	server *app.Server
}

func NewServer(server *app.Server) *Server {
	return &Server{server: server}
}

var channel = make(chan *string, 4)

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
)

func (s *Server) Init() {
	SubRoutineCustomers := s.server.Mux.PathPrefix("/api/customer").Subrouter()
	SubRoutineCustomers.HandleFunc("/save", s.HandleRegister).Methods(POST)
	SubRoutineCustomers.HandleFunc("/login", s.HandleLogin).Methods(GET)
}
