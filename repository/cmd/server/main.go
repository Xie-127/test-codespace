package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"userapi/repository"
)

type apiServer struct {
    repo repository.UserRepository
}

func (s *apiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch {
    case r.Method == http.MethodGet && r.URL.Path == "/users":
        s.handleListUsers(w)
    case r.Method == http.MethodPost && r.URL.Path == "/users":
        s.handleCreateUser(w, r)
    case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/users/"):
        s.handleGetUser(w, r)
    default:
        http.NotFound(w, r)
    }
}

func (s *apiServer) handleListUsers(w http.ResponseWriter) {
    users, err := s.repo.FindAll()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    writeJSON(w, users, http.StatusOK)
}

func (s *apiServer) handleCreateUser(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }
    if input.Name == "" || input.Email == "" {
        http.Error(w, "name and email are required", http.StatusBadRequest)
        return
    }

    user, err := s.repo.Create(repository.User{Name: input.Name, Email: input.Email})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    writeJSON(w, user, http.StatusCreated)
}

func (s *apiServer) handleGetUser(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/users/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.NotFound(w, r)
        return
    }
    user, err := s.repo.FindByID(id)
    if err != nil {
        if err == repository.ErrUserNotFound {
            http.NotFound(w, r)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    writeJSON(w, user, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, v interface{}, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(v)
}

func main() {
    port := flag.String("port", "8081", "server port")
    repoType := flag.String("repo", "memory", "repository type: memory or mysql")
    dsn := flag.String("dsn", "", "MySQL DSN for mysql repository")
    flag.Parse()

    var repo repository.UserRepository
    if *repoType == "mysql" {
        if *dsn == "" {
            panic("mysql DSN is required when repo=mysql")
        }
        gormRepo, err := repository.NewGormUserRepository(*dsn)
        if err != nil {
            panic(err)
        }
        repo = gormRepo
    } else {
        repo = repository.NewInMemoryUserRepository()
    }

    server := &apiServer{repo: repo}
    fmt.Printf("Service listening on :%s using repo=%s\n", *port, *repoType)
    if err := http.ListenAndServe(":"+*port, server); err != nil {
        panic(err)
    }
}
