package models

import (
	"time"
)

type ApplyVacation struct {
	ID                     uint         `gorm:"primaryKey"`
	MemberID               uint         `gorm:"index"`
	Member                 Member       `gorm:"foreignKey:MemberID"`
	VacationPlanID         uint         `gorm:"index"`
	VacationPlan           VacationPlan `gorm:"foreignKey:VacationPlanID"`
	VacationTypeID         uint         `gorm:"index"`
	VacationType           VacationType `gorm:"foreignKey:VacationTypeID"`
	StartDate              time.Time
	EndDate                time.Time
	HalfFirst              bool
	HalfLast               bool
	VacationProcessStateID uint                 `gorm:"index"`
	VacationProcessState   VacationProcessState `gorm:"foreignKey:VacationProcessStateID"`
	VacationCancelStateID  uint                 `gorm:"index"`
	VacationCancelState    VacationCancelState  `gorm:"foreignKey:VacationCancelStateID"`
}
