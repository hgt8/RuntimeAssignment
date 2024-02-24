package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type APIServer struct {
	ListenAddress string
	store         Storage
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type ApiError struct {
	Error string
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			err := WriteJson(w, http.StatusBadRequest, err.Error())
			if err != nil {
				return
			}
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

	//router.HandleFunc("/policies", makeHTTPHandleFunc(s.handlePoliciesEndpoints))
	router.HandleFunc("/policies", makeHTTPHandleFunc(s.handleCreatePolicy)).Methods("POST")
	router.HandleFunc("/policies/{id}", makeHTTPHandleFunc(s.handleGetPolicy)).Methods("GET")
	router.HandleFunc("/policies/{id}", makeHTTPHandleFunc(s.handleUpdatePolicy)).Methods("PUT")
	router.HandleFunc("/policies/{id}", makeHTTPHandleFunc(s.handleDeletePolicy)).Methods("DELETE")

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
		return s.handleDeletePolicy(w, r)
	}

	return fmt.Errorf("method not supported %s", r.Method)
}

func (s *APIServer) handleCreatePolicy(w http.ResponseWriter, r *http.Request) error {
	createPolicyRequest := &CreatePolicyRequest{}
	if err := json.NewDecoder(r.Body).Decode(&createPolicyRequest); err != nil {
		return err
	}
	policy := CreatePolicyRequest{
		createPolicyRequest.PolicyName,
		createPolicyRequest.Author,
		createPolicyRequest.ControlData,
	}
	if err := s.store.CreatePolicy(&policy); err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, &policy)
}

func (s *APIServer) handleGetPolicy(w http.ResponseWriter, r *http.Request) error {
	idString := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID, Error in Conversion", http.StatusBadRequest)
		return err
	}
	policy, err := s.store.GetPolicy(id)
	if err != nil {
		http.Error(w, "Policy not found", http.StatusNotFound)
		return err
	}

	// Return the policy as JSON
	return WriteJson(w, http.StatusOK, &policy)
}

func (s *APIServer) handleUpdatePolicy(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return err
	}

	var updateReq UpdatePolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := validateUpdatePolicyRequest(&updateReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := s.store.UpdatePolicy(id, &updateReq); err != nil {
		http.Error(w, "Failed to update policy", http.StatusInternalServerError)
		return err
	}

	return WriteJson(w, http.StatusOK, map[string]string{"message": "Policy updated successfully"})
}

func validateUpdatePolicyRequest(policy *UpdatePolicyRequest) error {
	if policy.PolicyName == "" {
		return errors.New("policyName cannot be empty")
	}
	if policy.Author == "" {
		return errors.New("author cannot be empty")
	}
	// Add any additional validation rules here
	return nil
}

func (s *APIServer) handleDeletePolicy(w http.ResponseWriter, r *http.Request) error {
	idString := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID, Error in Conversion", http.StatusBadRequest)
		return err
	}

	err = s.store.DeletePolicy(id)
	if err != nil {
		if err.Error() == "no rows affected" {
			http.Error(w, "Policy not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error deleting policy", http.StatusInternalServerError)
		}
		return err
	}
	message := fmt.Sprintf("Deleted policy at id: %d", id)
	return WriteJson(w, http.StatusOK, message)
}
