package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/user", makeHTTPHandleFunc(s.handleUser))
	router.HandleFunc("/task", withJWTAuth(makeHTTPHandleFunc(s.handleTask), s.store))
	router.HandleFunc("/task/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleTaskById), s.store))

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

// 498081
func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	user, err := s.store.GetUserByEmail(req.Email)
	if err != nil {
		return fmt.Errorf("there is no user with given email")
	}

	if !user.ValidPassword(req.Password) {
		return fmt.Errorf("not authenticated")
	}

	token, err := createJWT(user)
	if err != nil {
		return fmt.Errorf("not authenticated")
	}

	resp := LoginResponse{
		Token: token,
	}

	return WriteJSON(w, http.StatusOK, resp)
}

func (s *APIServer) handleUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetUsers(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateUser(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleTask(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetTasks(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateTask(w, r)
	}
	// if r.Method == "PUT" {
	// 	return s.handleUpdateTask(w, r)
	// }

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := s.store.GetUsers()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, users)
}

func (s *APIServer) handleGetTasks(w http.ResponseWriter, r *http.Request) error {
	tasks, err := s.store.GetTasks()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, tasks)
}

func (s *APIServer) handleTaskById(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}

		task, err := s.store.GetTaskById(id)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, task)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteTask(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	req := new(RegisterRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	user, err := NewUser(req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		return err
	}
	if err := s.store.CreateUser(user); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleCreateTask(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateTaskRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	task, err := NewTask(req.Title, req.Description)
	if err != nil {
		return err
	}
	if err := s.store.CreateTask(task); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, task)
}

func (s *APIServer) handleDeleteTask(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteTask(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func createJWT(user *User) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"email":     user.Email,
	}

	secret := getSigningKey()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func getSigningKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable not found - aborting")
	}
	return []byte(secret)
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}

func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")

		splitToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		auth := splitToken[1]
		token, err := validateJWT(auth)
		claims := token.Claims.(jwt.MapClaims)
		email := claims["email"].(string)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}

		user, err := s.GetUserByEmail(email)
		if err != nil {
			permissionDenied(w)
			return
		}

		if user.Email != claims["email"] {
			permissionDenied(w)
		}

		if err != nil {
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "invalid token"})
			return
		}

		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := getSigningKey()
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
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

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["Authorization"]
	fmt.Println(idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}
