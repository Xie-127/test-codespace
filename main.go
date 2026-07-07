package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "strings"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type app struct {
    users map[int]User
    nextID int
}

func newApp() *app {
    return &app{users: make(map[int]User), nextID: 1}
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch {
    case r.Method == http.MethodGet && r.URL.Path == "/users":
        a.handleListUsers(w, r)
    case r.Method == http.MethodPost && r.URL.Path == "/users":
        a.handleCreateUser(w, r)
    case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/users/"):
        a.handleGetUser(w, r)
    default:
        http.NotFound(w, r)
    }
}

func (a *app) handleListUsers(w http.ResponseWriter, r *http.Request) {
    users := make([]User, 0, len(a.users))
    for _, user := range a.users {
        users = append(users, user)
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(users)
}

func (a *app) handleCreateUser(w http.ResponseWriter, r *http.Request) {
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

    user := User{ID: a.nextID, Name: input.Name, Email: input.Email}
    a.users[user.ID] = user
    a.nextID++

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(w).Encode(user)
}

func (a *app) handleGetUser(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/users/")
    if idStr == "" || strings.Contains(idStr, "/") {
        http.NotFound(w, r)
        return
    }

    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.NotFound(w, r)
        return
    }

    user, ok := a.users[id]
    if !ok {
        http.NotFound(w, r)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(user)
}

func main() {
    app := newApp()
    fmt.Println("Server listening on :8080")
    if err := http.ListenAndServe(":8080", app); err != nil {
        panic(err)
    }
}
