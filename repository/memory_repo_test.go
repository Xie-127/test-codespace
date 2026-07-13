package repository

import "testing"

func TestInMemoryUserRepositoryCRUD(t *testing.T) {
    repo := NewInMemoryUserRepository()

    created, err := repo.Create(User{Name: "Alice", Email: "alice@example.com"})
    if err != nil {
        t.Fatalf("Create failed: %v", err)
    }
    if created.ID == 0 {
        t.Fatal("expected created user to have an ID")
    }

    found, err := repo.FindByID(created.ID)
    if err != nil {
        t.Fatalf("FindByID failed: %v", err)
    }
    if found.Name != created.Name || found.Email != created.Email {
        t.Fatalf("expected found user to match created user, got %+v", found)
    }

    updated, err := repo.Update(User{ID: created.ID, Name: "Alice Updated", Email: "alice@example.com"})
    if err != nil {
        t.Fatalf("Update failed: %v", err)
    }
    if updated.Name != "Alice Updated" {
        t.Fatalf("expected updated name, got %q", updated.Name)
    }

    users, err := repo.FindAll()
    if err != nil {
        t.Fatalf("FindAll failed: %v", err)
    }
    if len(users) != 1 {
        t.Fatalf("expected 1 user, got %d", len(users))
    }

    if err := repo.Delete(created.ID); err != nil {
        t.Fatalf("Delete failed: %v", err)
    }

    _, err = repo.FindByID(created.ID)
    if err != ErrUserNotFound {
        t.Fatalf("expected ErrUserNotFound after delete, got %v", err)
    }
}
