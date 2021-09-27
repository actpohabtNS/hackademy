package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
)

type User struct {
	Email          string
	PasswordDigest string
	FavoriteCake   string
}
type UserRepository interface {
	Add(string, User) error
	Get(string) (User, error)
	Update(string, User) error
	Delete(string) (User, error)
}

type UserService struct {
	repository UserRepository
}
type UserRegisterParams struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	FavoriteCake string `json:"favorite_cake"`
}

func validateRegisterParams(p *UserRegisterParams) error {
	// 1. Email is valid
	if _, err := mail.ParseAddress(p.Email); err != nil {
		return errors.New("must provide an email")
	}

	// 2. Password at least 8 symbols
	if len(p.Password) < 8 {
		return errors.New("password must be at least 8 symbols")
	}

	// 3. Favorite cake not empty
	if len(p.FavoriteCake) < 1 {
		return errors.New("favourite cake can't be empty")
	}

	// 4. Favorite cake only alphabetic
	for _, charVariable := range p.FavoriteCake {
		if (charVariable < 'a' || charVariable > 'z') && (charVariable < 'A' || charVariable > 'Z') {
			return errors.New("favourite cake must contain only alphabetic characters")
		}
	}

	return nil
}
func (u *UserService) Register(w http.ResponseWriter, r *http.Request) {
	params := &UserRegisterParams{}
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		handleError(errors.New("could not read params"), w)
		return
	}
	if err := validateRegisterParams(params); err != nil {
		handleError(err, w)
		return
	}
	passwordDigest := md5.New().Sum([]byte(params.Password))
	newUser := User{
		Email:          params.Email,
		PasswordDigest: string(passwordDigest),
		FavoriteCake:   params.FavoriteCake,
	}
	err = u.repository.Add(params.Email, newUser)
	if err != nil {
		handleError(err, w)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("registered"))
}

func handleError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	_, _ = w.Write([]byte(err.Error()))
}
