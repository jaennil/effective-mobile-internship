package service

import (
	"sync"
)

type UserStorage struct {
    mutex sync.RWMutex
	Users  map[int32]*User
	nextID int32
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
        Users:  make(map[int32]*User),
		nextID: 1,
	}
}

func (s *UserStorage) Save(user *User) error {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    user.Id = s.nextID
    s.Users[s.nextID] = user.Clone()
    s.nextID++
    return nil
}

func (s *UserStorage) Find(username string) (*User, error) {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    for _, u := range s.Users {
        if u.Username == username {
            return u.Clone(), nil
        }
    }

    return nil, nil
}

func (s *UserStorage) Update(user *User) error {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    s.Users[user.Id] = user

    return nil
}
