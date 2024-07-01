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

func MapMemberToDTO(member *models.Member) MemberResponse {
	if member == nil {
		return MemberResponse{}
	}
	return MemberResponse{
		ID:       member.ID,
		Name:     member.Name,
		Email:    member.Email,
		HireDate: member.HireDate,
		IsActive: member.IsActive,
	}
}

func MapMembersToDTO(members []*models.Member) []*MemberResponse {
	var memberDTOs []*MemberResponse
	if members == nil {
		return memberDTOs
	}
	for _, member := range members {
		memberDTO := MapMemberToDTO(member)
		memberDTOs = append(memberDTOs, &memberDTO)
	}
	return memberDTOs
}
