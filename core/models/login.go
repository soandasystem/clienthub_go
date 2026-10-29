package models

type LoginRequest struct {
	LoginType  string `json:"login_type" binding:"required"` // "user", "course", "access_code"
	Username   string `json:"username,omitempty"`            // Usado para "user"
	Password   string `json:"password,omitempty"`            // Usado para "user" y "course"
	Rutapod    string `json:"rutapod,omitempty"`             // Usado para "course"
	AccessCode string `json:"access_code,omitempty"`         // Usado para "access_code"
}
