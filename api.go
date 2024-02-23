package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type APIServer struct {
	ListenAddress string
	store         Storage
}

func WriteJson(w http.ResponseWriter, status int, message string) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(message)
}

type ApiError struct {
	Error string
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, err.Error())
		}
	}
}

func Server(address string, store Storage) *APIServer {
	return &APIServer{
		ListenAddress: address,
		store:         store,
	}
}

func (s *APIServer) Run() {
	log.Println("Initiating Server on", s.ListenAddress)
	router := mux.NewRouter()

	router.HandleFunc("/policies", makeHTTPHandleFunc(s.handlePoliciesEndpoints))

	router.HandleFunc("/policies/{id}", makeHTTPHandleFunc(s.handleGetPolicy))

	err := http.ListenAndServe(s.ListenAddress, router)
	if err != nil {
		//log error
	}
}

func (s *APIServer) handlePoliciesEndpoints(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetPolicy(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreatePolicy(w, r)
	}
	if r.Method == "PATCH" {
		return s.handleUpdatePolicy(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleRemovePolicy(w, r)
	}

	return fmt.Errorf("method not supported %s", r.Method)
}

func (s *APIServer) handleCreatePolicy(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleGetPolicy(w http.ResponseWriter, r *http.Request) error {
	//id := mux.Vars(r)["id"]

	//get from DB

	return nil
}

func (s *APIServer) handleUpdatePolicy(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleRemovePolicy(w http.ResponseWriter, r *http.Request) error {
	return nil
}
