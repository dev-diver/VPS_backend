package dto

import "time"

// CreateVacationPlanRequest DTO for creating a vacation plan
type CreateVacationPlanRequest struct {
	Vacations       []VacationRequest `json:"vacations"`
	Approver1ID     uint              `json:"approver_1"`
	ApproverFinalID uint              `json:"approver_final"`
}

// VacationRequest DTO for vacation details
type VacationRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"` //TODO: start_date보다 크게
	HalfFirst bool      `json:"half_first"`
	HalfLast  bool      `json:"half_last"`
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
	Status    string    `json:"status"`
}

// ApproveVacationPlanRequest DTO for approving a vacation plan
type ApproveVacationPlanRequest struct {
	Approver1ID     uint `json:"approver_1"`
	ApproverFinalID uint `json:"approver_final"`
}

// VacationPlanResponse DTO for vacation plan response
type VacationPlanResponse struct {
	ID           uint                    `json:"id"`
	MemberID     uint                    `json:"member_id"`
	ApplyDate    time.Time               `json:"apply_date"`
	ApproveDate  *time.Time              `json:"approve_date"`
	Vacations    []ApplyVacationResponse `json:"vacations"`
	ProcessState string                  `json:"process_state"`
	CancelState  string                  `json:"cancel_state"`
}

// ApplyVacationResponse DTO for vacation response
type ApplyVacationResponse struct {
	ID           uint      `json:"id"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	HalfFirst    bool      `json:"half_first"`
	HalfLast     bool      `json:"half_last"`
	Status       string    `json:"status"`
	CancelStatus string    `json:"cancel_status"`
}
