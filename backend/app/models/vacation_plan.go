package models

import "time"

type VacationPlan struct {
	ID                     uint   `gorm:"primaryKey"`
	MemberID               uint   `gorm:"index"`
	Member                 Member `gorm:"foreignKey:MemberID"`
	Approver1ID            uint   `gorm:"foreignKey:MemberID"`
	ApproverFinalID        uint   `gorm:"foreignKey:MemberID"`
	ApplyDate              time.Time
	VacationProcessStateID uint                 `gorm:"index"`
	VacationProcessState   VacationProcessState `gorm:"foreignKey:VacationProcessStateID"`
	VacationCancelStateID  uint                 `gorm:"index"`
	VacationCancelState    VacationCancelState  `gorm:"foreignKey:VacationCancelStateID"`
	ApplyVacations         []ApplyVacation      `gorm:"foreignKey:VacationPlanID"`
}
