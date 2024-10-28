package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
    Id int32
    Username string
    PasswordHash string
}

func NewUser(username string, password string) (*User, error) {
    passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("cannot hash password: %v", err)
    }

    user := &User {
        Username: username,
        PasswordHash: string(passwordHash),
    }

    return user, nil
}

func (user *User) IsPasswordCorrect(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
    return err == nil
}

func (user *User) Clone() *User {
    return &User {
        Id: user.Id,
        Username: user.Username,
        PasswordHash: user.PasswordHash,
    }
}

func (user *User) String() string {
    return fmt.Sprintf("%d %s %s", user.Id, user.Username, user.PasswordHash)
}
