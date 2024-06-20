package dto

import "time"

type CreateMemberDTO struct {
	Name     string    `json:"name" validate:"required"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required"`
	HireDate time.Time `json:"hire_date" validate:"required"`
}

type MemberDTO struct {
	ID       uint      `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	HireDate time.Time `json:"hire_date"`
	IsActive bool      `json:"is_active"`
}
