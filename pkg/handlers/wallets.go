package handlers

import (
	"encoding/json"
	"github.com/RakhimovAns/Alif/pkg/service"
	"github.com/RakhimovAns/Alif/types"
	"log"
	"net/http"
	"strconv"
)

func (s *Server) HandleCreateWallet(writer http.ResponseWriter, request *http.Request) {
	id := *<-channel
	err := s.server.WalletService.CreateWallet(request.Context(), id)
	if err == types.ErrNotIdentified {
		log.Println(err)
		http.Error(writer, "Not Identified", http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(writer, "Creating has been failed", http.StatusInternalServerError)
		return
	}
	_, err = writer.Write([]byte("Was created successfully"))
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) HandleCheckWallet(writer http.ResponseWriter, request *http.Request) {

	id := *<-channel
	err := s.server.WalletService.CheckWallet(request.Context(), id)
	if err == types.ErrNoSuch {
		http.Error(writer, "Not Found", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	writer.Write([]byte("Was founded successfully"))
}
func (s *Server) HandleDepositWallet(writer http.ResponseWriter, request *http.Request) {
	id := *<-channel
	SumParam := service.RequestBod.OtherData["sum"]
	sum, _ := strconv.ParseInt(SumParam, 10, 64)
	err := s.server.WalletService.DepositWallet(request.Context(), id, sum)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	writer.Write([]byte("Was deposited successfully"))
}
func (s *Server) HandleGetBalance(writer http.ResponseWriter, request *http.Request) {
	id := *<-channel
	balance, err := s.server.WalletService.GetBalance(request.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(balance)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}
func (s *Server) HandleGetActions(writer http.ResponseWriter, request *http.Request) {
	id := *<-channel
	item, err := s.server.WalletService.GetActions(request.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
