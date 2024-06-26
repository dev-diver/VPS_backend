package models

import "time"

type VacationPlan struct {
	ID                     uint   `gorm:"primaryKey"`
	MemberID               uint   `gorm:"index"`
	Member                 Member `gorm:"foreignKey:MemberID"`
	Approver1ID            uint   `gorm:"index"`
	Approver1              Member `gorm:"foreignKey:Approver1ID"`
	ApproverFinalID        uint   `gorm:"index"`
	ApproverFinal          Member `gorm:"foreignKey:ApproverFinalID"`
	ApplyDate              time.Time
	VacationProcessStateID uint                 `gorm:"index"`
	VacationProcessState   VacationProcessState `gorm:"foreignKey:VacationProcessStateID"`
	VacationCancelStateID  uint                 `gorm:"index"`
	VacationCancelState    VacationCancelState  `gorm:"foreignKey:VacationCancelStateID"`
	ApplyVacations         []ApplyVacation      `gorm:"foreignKey:VacationPlanID"`
}
