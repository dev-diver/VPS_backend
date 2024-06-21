package dto

import "time"

// CreateVacationPlanRequest DTO for creating a vacation plan
type CreateVacationPlanRequest struct {
	Vacations       []VacationRequest `json:"vacations" validate:"required"`
	Approver1ID     uint              `json:"approver_1" validate:"required"`
	ApproverFinalID uint              `json:"approver_final" validate:"required"`
}

// VacationRequest DTO for vacation details
type VacationRequest struct {
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required,gtefield=StartDate"`
	HalfFirst bool      `json:"half_first" validate:"required"`
	HalfLast  bool      `json:"half_last" validate:"required"`
}

// EditVacationPlanRequest DTO for editing a vacation plan
type EditVacationPlanRequest struct {
	Vacations []VacationEditRequest `json:"vacations"`
}

// VacationEditRequest DTO for editing vacation details
type VacationEditRequest struct {
	ID        uint      `json:"id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	HalfFirst bool      `json:"half_first"`
	HalfLast  bool      `json:"half_last"`
	Status    string    `json:"process_state"`
}

// ApproveVacationPlanRequest DTO for approving a vacation plan
type ApproveVacationPlanRequest struct {
	ApprovalState uint `json:"approval_state" validate:"required"`
	MemberID      uint `json:"member_id" validate:"required"`
}

// VacationPlanResponse DTO for vacation plan response
type VacationPlanResponse struct {
	ID           uint                    `json:"id"`
	MemberID     uint                    `json:"member_id"`
	MemberName   string                  `json:"member_name"`
	ApplyDate    time.Time               `json:"apply_date"`
	ApproveDate  *time.Time              `json:"approve_date"`
	Vacations    []ApplyVacationResponse `json:"vacations"`
	ProcessState uint                    `json:"process_state"`
	CancelState  uint                    `json:"cancel_state"`
}

// ApplyVacationResponse DTO for vacation response
type ApplyVacationResponse struct {
	ID           uint      `json:"id"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	HalfFirst    bool      `json:"half_first"`
	HalfLast     bool      `json:"half_last"`
	Status       uint      `json:"process_state"`
	CancelStatus uint      `json:"cancel_state"`
}

type ApplyVacationCardResponse struct {
	ID           uint      `json:"id"`
	MemberID     uint      `json:"member_id"`
	MemberName   string    `json:"member_name"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	HalfFirst    bool      `json:"half_first"`
	HalfLast     bool      `json:"half_last"`
	Status       uint      `json:"process_state"`
	CancelStatus uint      `json:"cancel_state"`
}
