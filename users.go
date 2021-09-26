package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type UserService struct {
	repository UserRepository
}

type UserRegisterParams struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	FavoriteCake string `json:"favorite_cake"`
}

func ValidateRegisterParams(p *UserRegisterParams) error {
	// check requirements
	return nil
}

func (u *UserService) Register(w http.ResponseWriter, r *http.Request) {
	params := &UserRegisterParams{}
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		handleError(errors.New("could not read params"), w)
		return
	}

	if err := ValidateRegisterParams(params); err != nil {
		handleError(err, w)
		return
	}
	passwordDigest := md5.New().Sum([]byte(params.Password))
	newUser := User{
		Email:          params.Email,
		PasswordDigest: string(passwordDigest),
		FavoriteCake:   params.FavoriteCake,
	}

	fmt.Println("cake: ", newUser.FavoriteCake)
	fmt.Println("email: ", newUser.Email)

	err = u.repository.Add(params.Email, newUser)
	if err != nil {
		handleError(err, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("registered"))
	UserRegisteredCounter.Inc()
}

func handleError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write([]byte(err.Error()))
}

type User struct {
	Email          string
	PasswordDigest string
	FavoriteCake   string
	Role		string
}

type UserRepository interface {
	Add(string, User) error
	Get(string) (User, error)
	Update(string, User) error
	Delete(string) (User, error)
}
