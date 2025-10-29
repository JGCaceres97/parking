package domain

import "time"

type Role = string

const (
	RoleAdmin  Role = "admin"
	RoleCommon Role = "common"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Role      Role      `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
