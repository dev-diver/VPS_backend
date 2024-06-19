package models

import (
	"time"

	"cywell.com/vacation-promotion/app/models"
)

type ConsumeVacation struct {
	ID             uint                `gorm:"primaryKey"`
	MemberID       uint                `gorm:"index"`
	Member         Member              `gorm:"foreignKey:MemberID"`
	VacationPlanID uint                `gorm:"index"`
	VacationPlan   VacationPlan        `gorm:"foreignKey:VacationPlanID"`
	VacationTypeID uint                `gorm:"index"`
	VacationType   models.VacationType `gorm:"foreignKey:VacationTypeID"`
	StartDate      time.Time
	EndDate        time.Time
	HalfFirst      bool
	HalfLast       bool
}
