package main

import "time"

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
