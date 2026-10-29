package models

import "time"

type Users struct {
	ID           string    `json:"_id,omitempty"`
	Username     string    `json:"username"`
	FullName     string    `json:"full_name"`
	PasswordHash string    `json:"password_hash,omitempty"`
	Email        string    `json:"email"`
	RolesId      int64     `json:"roles_id"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

// Resp  response struct
type UsersResp struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	FullName     string    `json:"full_name"`
	PasswordHash string    `json:"password_hash,omitempty"`
	Email        string    `json:"email"`
	RolesId      int64     `json:"roles_id"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (UsersResp) TableName() string {
	return "users" // Nombre de la tabla en la base de datos
}

// Create---Req  request struct
type CreateUsersReq struct {
	ID           string    `gorm:"primaryKey;autoIncrement"`
	Username     string    `json:"username"`
	FullName     string    `json:"full_name"`
	PasswordHash string    `json:"password_hash,omitempty"`
	Email        string    `json:"email"`
	RolesId      int64     `json:"roles_id"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (CreateUsersReq) TableName() string {
	return "users" // Nombre de la tabla en la base de datos
}

type UpdateUsersReq struct {
	ID           string     `json:"id"`
	Username     *string    `json:"username"`
	FullName     *string    `json:"full_name"`
	PasswordHash *string    `json:"password_hash,omitempty"`
	Email        *string    `json:"email"`
	RolesId      *int64     `json:"roles_id"`
	Active       *bool      `json:"active"`
	CreatedAt    *time.Time `gorm:"autoCreateTime"`
}

func (UpdateUsersReq) TableName() string {
	return "users" // Nombre de la tabla en la base de datos
}

type UsersInf struct {
	ID           string      `json:"id"`
	Username     string      `json:"username"`
	FullName     string      `json:"full_name"`
	PasswordHash string      `json:"password_hash,omitempty"`
	Email        string      `json:"email"`
	RolesId      int64       `json:"roles_id"`
	Rol          RolesReport `json:"rol" gorm:"foreignKey:RolesId;references:ID"`
	Active       bool        `json:"active"`
	CreatedAt    time.Time   `gorm:"autoCreateTime"`
}

func (UsersInf) TableName() string {
	return "users" // Nombre de la tabla en la base de datos
}

type UsersListResponse struct {
	Items      []UsersInf `json:"items"`
	TotalCount int64      `json:"totalCount"`
}

type UsersReport struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	FullName     string    `json:"full_name"`
	PasswordHash string    `json:"password_hash,omitempty"`
	Email        string    `json:"email"`
	RolesId      int64     `json:"roles_id"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (UsersReport) TableName() string {
	return "users" // Nombre de la tabla en la base de datos
}
