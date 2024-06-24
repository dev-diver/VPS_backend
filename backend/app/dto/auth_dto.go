package dto

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Member    MemberResponse `json:"member"`
	CompanyID uint           `json:"company_id"`
	GroupIDs  []uint         `json:"group_ids"`
}
