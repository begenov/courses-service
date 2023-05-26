package domain

import "time"

type Student struct {
	ID           int       `json:"id"`
	Email        string    `json:"email" binding:"required,email,max=64"`
	Name         string    `json:"name" binding:"required,min=3,max=64"`
	Password     string    `json:"password" binding:"required,min=8,max=64"`
	GPA          float64   `json:"gpa" binding:"required"`
	Courses      []string  `json:"courses"`
	RefreshToken string    `json:"-"`
	ExpiresAt    time.Time `json:"-"`
}

type Response struct {
	Students []Student `json:"students"`
}
