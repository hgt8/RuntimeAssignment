package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type APIServer struct {
	ListenAddress     string
	store             Storage
	activeConnections map[*websocket.Conn]bool
	connMutex         sync.Mutex
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
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
		ListenAddress:     address,
		store:             store,
		activeConnections: make(map[*websocket.Conn]bool),
	}
}

func (s *APIServer) Run() {
	log.Println("Initiating Server on", s.ListenAddress)
	router := mux.NewRouter()

	router.HandleFunc("/ws", s.wsHandler)
	router.HandleFunc("/policies", makeHTTPHandleFunc(s.handleGetAllPolicies)).Methods("GET")
	router.HandleFunc("/policies", makeHTTPHandleFunc(s.handleCreatePolicy)).Methods("POST")
	router.HandleFunc("/policies/{id}", makeHTTPHandleFunc(s.handleGetPolicy)).Methods("GET")
	router.HandleFunc("/policies/{id}", makeHTTPHandleFunc(s.handleUpdatePolicy)).Methods("PUT")
	router.HandleFunc("/policies/{id}", makeHTTPHandleFunc(s.handleDeletePolicy)).Methods("DELETE")

	err := http.ListenAndServe(s.ListenAddress, router)
	if err != nil {
		log.Fatal("Server failed to start")
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *APIServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Could not upgrade ws")
		return
	}
	defer conn.Close()

	policies, err := s.getAllPolicies()
	if err != nil {
		log.Fatal("Could not get all policies")
		return
	}
	policiesJson, err := json.Marshal(policies)
	if err != nil {
		log.Fatal("Could not serialize json")
		return
	}
	conn.WriteMessage(websocket.TextMessage, policiesJson)

	s.connMutex.Lock()
	s.activeConnections[conn] = true
	s.connMutex.Unlock()

	defer func() {
		s.connMutex.Lock()
		delete(s.activeConnections, conn)
		s.connMutex.Unlock()
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read error:", err)
			break
		}
		fmt.Printf("recv: %s\n", p)
		if err := conn.WriteMessage(messageType, p); err != nil {
			fmt.Println("write error:", err)
			break
		}
	}
}

func (s *APIServer) notifyClients(message string) {
	s.connMutex.Lock()
	defer s.connMutex.Unlock()

	serializedMessage, err := json.Marshal(message)
	if err != nil {
		log.Fatal("Could not serialize json")
		return
	}

	for conn := range s.activeConnections {
		if err := conn.WriteMessage(websocket.TextMessage, serializedMessage); err != nil {
			delete(s.activeConnections, conn)
		}
	}
}

func (s *APIServer) getAllPolicies() ([]*FullPolicy, error) {
	return s.store.GetAllPolicies()
}

func (s *APIServer) handleGetAllPolicies(w http.ResponseWriter, r *http.Request) error {
	policies, err := s.getAllPolicies()
	if err != nil {
		http.Error(w, "Failed to retrieve all policies", http.StatusInternalServerError)
		return err
	}

	return WriteJson(w, http.StatusOK, policies)
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

	return WriteJson(w, http.StatusOK, &policy)
}

func (s *APIServer) handleCreatePolicy(w http.ResponseWriter, r *http.Request) error {
	createPolicyRequest := &CreatePolicyRequest{}
	if err := json.NewDecoder(r.Body).Decode(&createPolicyRequest); err != nil {
		return err
	}
	if err := validateCreatePolicyRequest(createPolicyRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := s.store.CreatePolicy(createPolicyRequest); err != nil {
		return err
	}
	s.notifyClients("Create Policy has been triggered")
	return WriteJson(w, http.StatusOK, &createPolicyRequest)
}

func validateCreatePolicyRequest(policy *CreatePolicyRequest) error {
	if policy.PolicyName == "" {
		return errors.New("policyName cannot be empty")
	}
	if policy.Author == "" {
		return errors.New("author cannot be empty")
	}
	if policy.ControlData == nil {
		return errors.New("controlData cannot be empty")
	}
	return nil
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
	s.notifyClients("Update Policy has been triggered")
	message := fmt.Sprintf("Policy at id: %d was updated successfully", id)
	return WriteJson(w, http.StatusOK, message)
}

func validateUpdatePolicyRequest(policy *UpdatePolicyRequest) error {
	if policy.PolicyName == "" {
		return errors.New("policyName cannot be empty")
	}
	if policy.Author == "" {
		return errors.New("author cannot be empty")
	}
	if policy.ControlData == nil {
		return errors.New("controlData cannot be empty")
	}
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
