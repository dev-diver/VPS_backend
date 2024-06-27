package dto

import (
	"time"

	"cywell.com/vacation-promotion/app/models"
)

type CreateVacationPlanRequest struct {
	Vacations []VacationRequest `json:"vacations" validate:"required"`
	Approvers []uint            `json:"approvers" validate:"required"`
}

type VacationRequest struct {
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required,gtefield=StartDate"`
	HalfFirst bool      `json:"half_first" validate:"required"`
	HalfLast  bool      `json:"half_last" validate:"required"`
}

type EditVacationPlanRequest struct {
	Vacations []VacationEditRequest `json:"vacations"`
}

type VacationEditRequest struct {
	ID           uint      `json:"id"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	HalfFirst    bool      `json:"half_first"`
	HalfLast     bool      `json:"half_last"`
	ProcessState string    `json:"process_state"`
}

type ApproveVacationPlanRequest struct {
	ApprovalState uint `json:"approval_state" validate:"required"`
	MemberID      uint `json:"member_id" validate:"required"`
}

type VacationPlanResponse struct {
	ID            uint                    `json:"id"`
	MemberID      uint                    `json:"member_id"`
	MemberName    string                  `json:"member_name"`
	ApplyDate     time.Time               `json:"apply_date"`
	ApproverOrder []models.ApproverOrder  `json:"approver_order"`
	Vacations     []ApplyVacationResponse `json:"vacations"`
	ApproveStage  uint                    `json:"approve_stage"`
	RejectState   bool                    `json:"reject_state"`
}

type ApplyVacationResponse struct {
	ID           uint      `json:"id"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	HalfFirst    bool      `json:"half_first"`
	HalfLast     bool      `json:"half_last"`
	ApproveStage uint      `json:"approve_stage"`
	RejectState  bool      `json:"reject_state"`
}

type ApplyVacationCardResponse struct {
	ID           uint      `json:"id"`
	MemberID     uint      `json:"member_id"`
	MemberName   string    `json:"member_name"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	HalfFirst    bool      `json:"half_first"`
	HalfLast     bool      `json:"half_last"`
	ApproveStage uint      `json:"approve_stage"`
	RejectState  bool      `json:"reject_state"`
}

func MapApplyVacationToResponse(vacation models.ApplyVacation) ApplyVacationResponse {
	return ApplyVacationResponse{
		ID:           vacation.ID,
		StartDate:    vacation.StartDate,
		EndDate:      vacation.EndDate,
		HalfFirst:    vacation.HalfFirst,
		ApproveStage: vacation.ApproveStage,
		RejectState:  vacation.RejectState,
	}
}

func MapApplyVacationToCardResponse(vacation models.ApplyVacation) ApplyVacationCardResponse {
	return ApplyVacationCardResponse{
		ID:           vacation.ID,
		MemberID:     vacation.MemberID,
		MemberName:   vacation.Member.Name,
		StartDate:    vacation.StartDate,
		EndDate:      vacation.EndDate,
		HalfFirst:    vacation.HalfFirst,
		HalfLast:     vacation.HalfLast,
		ApproveStage: vacation.ApproveStage,
		RejectState:  vacation.RejectState,
	}
}

func MapVacationPlanToResponse(plan models.VacationPlan) VacationPlanResponse {
	return VacationPlanResponse{
		ID:            plan.ID,
		MemberID:      plan.MemberID,
		MemberName:    plan.Member.Name,
		ApplyDate:     plan.ApplyDate,
		ApproverOrder: plan.ApproverOrder,
		Vacations:     nil,
		ApproveStage:  plan.ApproveStage,
		RejectState:   plan.RejectState,
	}
}
