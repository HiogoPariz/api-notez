package api

import (
	"github.com/HiogoPariz/api-notez/internal/auth"
	"github.com/HiogoPariz/api-notez/internal/storage"
	"github.com/HiogoPariz/api-notez/internal/types"

	"encoding/json"
	"fmt"
	"log"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type (
	Store             = storage.Store
	CreatePageRequest = types.CreatePageRequest
)

type APIServer struct {
	listenAddr string
	store      Store
}

func NewAPIServer(listenAddr string, store Store) *APIServer {
	return &APIServer{
		listenAddr,
		store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/page", makeHTTPHandleFunc(s.handlePage))
	router.HandleFunc("/page/{id}", withJWTAuth(makeHTTPHandleFunc(s.handlePageByID)))

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handlePage(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetPage(w, r)
	}

	if r.Method == "POST" {
		return s.handleCreatePage(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handlePageByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetPageByID(w, r)
	}

	if r.Method == "POST" {
		return s.handleUpdatePageByID(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeletePageByID(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetPage(w http.ResponseWriter, r *http.Request) error {
	pages, err := s.store.GetPages()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, pages)
}

func (s *APIServer) handleGetPageByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getIDFromRequest(r)

	if err != nil {
		return err
	}

	page, err := s.store.GetPageByID(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, page)
}

func (s *APIServer) handleCreatePage(w http.ResponseWriter, r *http.Request) error {
	createPageReq := CreatePageRequest{}
	if err := json.NewDecoder(r.Body).Decode(&createPageReq); err != nil {
		return err
	}

	page := types.NewPage(createPageReq.Title, createPageReq.Content)

	if err := s.store.CreatePage(page); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, page)
}

func (s *APIServer) handleDeletePageByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getIDFromRequest(r)

	if err != nil {
		return err
	}

	if err := s.store.DeletePage(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleUpdatePageByID(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func withJWTAuth(handlerFun http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("JWT middleware")

		tokenString := r.Header.Get("x-jwt-token")
		_, err := auth.ValidateJWT(tokenString)
		if err != nil {
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "Forbidden"})
			return
		}

		handlerFun(w, r)
	}
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getIDFromRequest(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id %s", idStr)
	}

	return id, nil
}
