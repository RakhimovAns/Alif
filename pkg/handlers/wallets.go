package handlers

import (
	"encoding/json"
	"github.com/RakhimovAns/Alif/types"
	"log"
	"net/http"
)

func (s *Server) HandleCreateWallet(writer http.ResponseWriter, request *http.Request) {
	var item *types.Customer
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.server.WalletService.CreateWallet(request.Context(), item)
	if err == types.ErrNotIdentified {
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
