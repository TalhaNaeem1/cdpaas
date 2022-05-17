package models

type EmailTemplate struct {
	ToName  string `json:"to_name" binding:"required"`
	ToEmail string `json:"to_email" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}