package model

import (
	"time"
)

const (
	RoleIDAdmin   = "05ed6375-0786-4e1a-bcb0-1533c837954d"
	RoleNameAdmin = "admin"
	RoleIDUser    = "ab9f2901-5aea-43b6-8f2b-7bf97dd30808"
	RoleNameUser  = "user"
	// see 02_up.sql
)

type UserRole struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	ValidFrom *time.Time `json:"validFrom,omitempty"`
	ValidTo   *time.Time `json:"validTo,omitempty"`
}
