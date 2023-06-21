package handlers

import (
	"github.com/RakhimovAns/Alif/cmd/app"
	"github.com/RakhimovAns/Alif/pkg/service"
)

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
	PUT    = "PUT"
	DELETE = "DELETE"
)

func (s *Server) Init() {
	SubRoutineCustomers := s.server.Mux.PathPrefix("/api/customer").Subrouter()
	SubRoutineCustomers.HandleFunc("/save", s.HandleRegister).Methods(POST)
	SubRoutineCustomers.HandleFunc("/login", s.HandleLogin).Methods(GET)

	SubRoutineWallets := s.server.Mux.PathPrefix("/api/wallets").Subrouter()
	SubRoutineWallets.Use(service.Auth(channel))
	SubRoutineWallets.HandleFunc("/create", s.HandleCreateWallet).Methods(POST)
	SubRoutineWallets.HandleFunc("/check", s.HandleCheckWallet).Methods(GET)
	SubRoutineWallets.HandleFunc("/deposit", s.HandleDepositWallet).Methods(PUT)
	SubRoutineWallets.HandleFunc("/balance", s.HandleGetBalance).Methods(GET)
	SubRoutineWallets.HandleFunc("/actions", s.HandleGetActions).Methods(GET)
}
