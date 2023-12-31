package handlers

import (
	"encoding/json"
	"github.com/RakhimovAns/Alif/pkg/service"
	"github.com/RakhimovAns/Alif/types"
	"net/http"
)

func (s *Server) HandleRegister(writer http.ResponseWriter, request *http.Request) {
	var item *types.Customer
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = s.server.CustomerService.Register(request.Context(), item)
	if err == types.ErrNotIdentified {
		http.Error(writer, "Not Identified", http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(writer, "Register has been failed", http.StatusInternalServerError)
		return
	}
	_, err = writer.Write([]byte("Was saved successfully"))
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) HandleLogin(writer http.ResponseWriter, request *http.Request) {
	requestBody := service.Parse(request)
	item := types.Customer{
		Login:    requestBody.CustomerData["login"],
		Password: requestBody.CustomerData["password"],
	}
	token, XDigest, err := s.server.CustomerService.Login(request.Context(), item.Login, item.Password, requestBody.AuthData)
	if err == types.ErrNoSuchUser {
		_, err = writer.Write([]byte("No account with this phone number"))
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else if err == types.ErrInvalidPassword {
		_, err = writer.Write([]byte("Passwords don't match"))
		if err != nil {
			http.Error(writer, http.StatusText(500), http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	} else if err == nil {
		_, err = writer.Write([]byte("You have been login successfully\nHere is your Token and X-digest\n"))
		if err != nil {
			http.Error(writer, http.StatusText(500), 500)
			return
		}
		data, err := json.Marshal(token)
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
		data, err = json.Marshal(XDigest)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		writer.Write([]byte("\n"))
		_, err = writer.Write(data)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	}
}
