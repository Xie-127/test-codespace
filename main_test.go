package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGetUsersReturnsEmptyInitially(t *testing.T) {
    app := newApp()
    req := httptest.NewRequest(http.MethodGet, "/users", nil)
    rr := httptest.NewRecorder()

    app.ServeHTTP(rr, req)

    if rr.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", rr.Code)
    }

    var users []User
    if err := json.NewDecoder(rr.Body).Decode(&users); err != nil {
        t.Fatalf("decode response: %v", err)
    }

    if len(users) != 0 {
        t.Fatalf("expected no users initially, got %d", len(users))
    }
}

func TestCreateAndGetUser(t *testing.T) {
    app := newApp()

    payload := []byte(`{"name":"Alice","email":"alice@example.com"}`)
    createReq := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(payload))
    createReq.Header.Set("Content-Type", "application/json")
    createRR := httptest.NewRecorder()

    app.ServeHTTP(createRR, createReq)

    if createRR.Code != http.StatusCreated {
        t.Fatalf("expected status 201, got %d", createRR.Code)
    }

    var created User
    if err := json.NewDecoder(createRR.Body).Decode(&created); err != nil {
        t.Fatalf("decode created user: %v", err)
    }

    if created.ID == 0 || created.Name != "Alice" || created.Email != "alice@example.com" {
        t.Fatalf("unexpected created user: %+v", created)
    }

    getReq := httptest.NewRequest(http.MethodGet, "/users/"+strconv.Itoa(created.ID), nil)
    getRR := httptest.NewRecorder()

    app.ServeHTTP(getRR, getReq)

    if getRR.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", getRR.Code)
    }

    var fetched User
    if err := json.NewDecoder(getRR.Body).Decode(&fetched); err != nil {
        t.Fatalf("decode fetched user: %v", err)
    }

    if fetched.ID != created.ID || fetched.Name != created.Name || fetched.Email != created.Email {
        t.Fatalf("unexpected fetched user: %+v", fetched)
    }
}
