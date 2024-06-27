package models

import "time"

type VacationPlan struct {
	ID             uint   `gorm:"primaryKey"`
	MemberID       uint   `gorm:"index"`
	Member         Member `gorm:"foreignKey:MemberID"`
	ApplyDate      time.Time
	ApproverOrder  []ApproverOrder `gorm:"foreignKey:VacationPlanID"`
	ApproveStage   uint            `gorm:"not null"`
	RejectState    bool            `gorm:"not null"`
	ApplyVacations []ApplyVacation `gorm:"foreignKey:VacationPlanID"`
}
