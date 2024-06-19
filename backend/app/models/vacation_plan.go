package models

import "time"

type VacationPlan struct {
	ID                     uint   `gorm:"primaryKey"`
	MemberID               uint   `gorm:"index"`
	Member                 Member `gorm:"foreignKey:MemberID"`
	Approver1ID            uint   `gorm:"index"`
	Approver2ID            uint   `gorm:"index"`
	FinalApproverID        uint   `gorm:"index"`
	ApproveDate            time.Time
	VacationProcessStateID uint                 `gorm:"index"`
	VacationProcessState   VacationProcessState `gorm:"foreignKey:VacationProcessStateID"`
	ConsumeVacations       []ConsumeVacation    `gorm:"foreignKey:VacationPlanID"`
}
