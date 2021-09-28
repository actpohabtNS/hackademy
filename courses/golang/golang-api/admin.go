package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Role int8

const (
	UserRole Role = iota
	AdminRole
	SuperAdminRole
)

type UserBanParams struct {
	Email  string `json:"email"`
	Reason string `json:"reason"`
}

type UserUnbanParams struct {
	Email string `json:"email"`
}

type BanHistoryQuery struct {
	Executor string
	Time     time.Time
	Reason   string
}

type BanHistory []BanHistoryQuery

func banHandler(w http.ResponseWriter, r *http.Request, executor User, users UserRepository) {
	params := &UserBanParams{}
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		handleUnprocError(errors.New("could not read params"), w)
		return
	}

	if err := validateEmail(params.Email); err != nil {
		handleUnprocError(err, w)
		return
	}

	user, getErr := users.Get(params.Email)
	if getErr != nil {
		handleUnprocError(err, w)
		return
	}

	if executor.Role <= user.Role {
		handleUnauthError(errors.New("permission denied"), w)
		return
	}

	banHistoryQuery := BanHistoryQuery{
		Executor: executor.Email,
		Time:     time.Now(),
		Reason:   params.Reason,
	}
	user.Banned = true
	user.BanHistory = append(user.BanHistory, banHistoryQuery)

	err = users.Update(user.Email, user)
	if err != nil {
		handleUnprocError(err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("user " + user.Email + " banned"))
}

func unbanHandler(w http.ResponseWriter, r *http.Request, executor User, users UserRepository) {
	params := &UserUnbanParams{}
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		handleUnprocError(errors.New("could not read params"), w)
		return
	}

	if err := validateEmail(params.Email); err != nil {
		handleUnprocError(err, w)
		return
	}

	user, getErr := users.Get(params.Email)
	if getErr != nil {
		handleUnprocError(err, w)
		return
	}

	if executor.Role <= user.Role {
		handleUnauthError(errors.New("permission denied"), w)
		return
	}

	banHistoryQuery := BanHistoryQuery{
		Executor: executor.Email,
		Time:     time.Now(),
		Reason:   "",
	}
	user.Banned = false
	user.BanHistory = append(user.BanHistory, banHistoryQuery)

	err = users.Update(user.Email, user)
	if err != nil {
		handleUnprocError(err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("user " + user.Email + " unbanned"))
}
