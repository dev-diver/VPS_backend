package dto

type OrganizeRequest struct {
	Name string `json:"name" validate:"required"`
}
