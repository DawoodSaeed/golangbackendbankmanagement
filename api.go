package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type APIServer struct {
	listenAddr string
}

type APIError struct {
	Error string
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
func makeHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHandlerFunc(s.handleAccount))
	http.ListenAndServe(s.listenAddr, router)
}

//ROUTE METHODS ############################
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccount(w, r)
		break
	case "POST":
		return s.handleCreateAccount(w, r)
		break
	case "DELETE":
		return s.handleCreateAccount(w, r)
		break
	default:
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	return nil
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	account := NewAccount("Dawood", "saeed")
	writeJSON(w, http.StatusOK, account)
	return nil
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
