package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type user struct {
	ID       int
	Username string `json:"username"`
	Age      int    `json:"age"`
}

type UserStorage struct {
	users  map[int]user
	nextID int
}

var userStorage = UserStorage{
	users:  make(map[int]user),
	nextID: 1,
}

func main() {
	http.HandleFunc("/users", handleUsers)
	http.HandleFunc("/users/", handleUserById)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var user user
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		user.ID = userStorage.nextID
		userStorage.nextID++
		userStorage.users[user.ID] = user

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	case http.MethodGet:
		var users []user
		for _, user := range userStorage.users {
			users = append(users, user)
		}

		json.NewEncoder(w).Encode(users)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleUserById(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		idStr := r.URL.Path[len("/users/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		if _, exists := userStorage.users[id]; !exists {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		delete(userStorage.users, id)

		w.WriteHeader(http.StatusNoContent)
	case http.MethodPut:
		idStr := r.URL.Path[len("/users/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		if _, exists := userStorage.users[id]; !exists {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		var updatedUser user
		if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
		}

		updatedUser.ID = id
		userStorage.users[id] = updatedUser

		json.NewEncoder(w).Encode(updatedUser)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
