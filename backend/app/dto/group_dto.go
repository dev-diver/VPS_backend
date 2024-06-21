package dto

import "cywell.com/vacation-promotion/app/models"

type CreateGroupRequest struct {
	Name     string `json:"name" validate:"required"`
	Color    string `json:"color" validate:"required"`
	Priority *int   `json:"priority"`
}

type GroupResponse struct {
	ID        uint            `json:"id"`
	CompanyID uint            `json:"company_id" validate:"required"`
	Name      string          `json:"name" validate:"required"`
	Color     string          `json:"color" validate:"required"`
	Priority  int             `json:"priority" validate:"required,min=1"`
	Members   []models.Member `json:"members" validate:"dive"`
}
