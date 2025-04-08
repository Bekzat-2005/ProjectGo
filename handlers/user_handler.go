package handlers

import (
	"encoding/json"
	"net/http"
	"projectGolang/models"
)

var users = []models.User{
	{ID: 1, Name: "John Doe"},
	{ID: 2, Name: "Jane Smith"},
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	user.ID = len(users) + 1
	users = append(users, user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
