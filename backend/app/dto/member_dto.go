package dto

import (
	"time"

	"cywell.com/vacation-promotion/app/models"
)

type CreateMemberRequest struct {
	Name     string    `json:"name" validate:"required"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required"`
	HireDate time.Time `json:"hire_date" validate:"required"`
}
type MemberResponse struct {
	ID       uint      `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	HireDate time.Time `json:"hire_date"`
	IsActive bool      `json:"is_active"`
}

func MapMembersToDTO(members []models.Member) []MemberResponse {
	var memberDTOs []MemberResponse
	for _, member := range members {
		memberDTO := MemberResponse{
			ID:       member.ID,
			Name:     member.Name,
			Email:    member.Email,
			HireDate: member.HireDate,
			IsActive: member.IsActive,
		}
		memberDTOs = append(memberDTOs, memberDTO)
	}
	return memberDTOs
}
