package dto

import "cywell.com/vacation-promotion/app/models"

type OrganizeRequest struct {
	Name string `json:"name" validate:"required"`
}

type OrganizeResponse struct {
	ID       uint               `json:"organize_id"`
	Name     string             `json:"organize_name"`
	ParentID *uint              `json:"parent_id"`
	Members  []MemberResponse   `json:"members,omitempty"`
	Children []OrganizeResponse `json:"children,omitempty"`
}

func MapOrganizeToResponse(organize models.Organize) OrganizeResponse {

	return OrganizeResponse{
		ID:       organize.ID,
		Name:     organize.Name,
		ParentID: organize.ParentID,
		Members:  MapMembersToDTO(organize.Members),
	}
}
