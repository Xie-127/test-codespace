package repository

import "sync"

type InMemoryUserRepository struct {
    mu     sync.RWMutex
    users  map[int]User
    nextID int
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
    return &InMemoryUserRepository{
        users:  make(map[int]User),
        nextID: 1,
    }
}

func (r *InMemoryUserRepository) Create(user User) (User, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    user.ID = r.nextID
    r.nextID++
    r.users[user.ID] = user
    return user, nil
}

func (r *InMemoryUserRepository) FindByID(id int) (User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    user, ok := r.users[id]
    if !ok {
        return User{}, ErrUserNotFound
    }
    return user, nil
}

func (r *InMemoryUserRepository) FindAll() ([]User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    users := make([]User, 0, len(r.users))
    for _, user := range r.users {
        users = append(users, user)
    }
    return users, nil
}

func (r *InMemoryUserRepository) Update(user User) (User, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, ok := r.users[user.ID]; !ok {
        return User{}, ErrUserNotFound
    }
    r.users[user.ID] = user
    return user, nil
}

func (r *InMemoryUserRepository) Delete(id int) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, ok := r.users[id]; !ok {
        return ErrUserNotFound
    }
    delete(r.users, id)
    return nil
}
