package models

// UserDetails is associated with authService

type UserInfo struct {
	ID          int         `json:"id"`
	Email       string      `json:"email"`
	WorkspaceID interface{} `json:"workspace_id"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
}

type Payload struct {
	UserInfo UserInfo `json:"user"`
}

type UserDetails struct {
	Success bool    `json:"success"`
	Payload Payload `json:"payload"`
	Errors  struct {
	} `json:"errors"`
	Description string `json:"description"`
}

type User struct {
	Type string `json:"type"`
}
