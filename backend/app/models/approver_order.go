package models

import (
	"time"
)

type ApproverOrder struct {
	ID             uint         `gorm:"primaryKey;autoIncrement"`
	VacationPlanID uint         `gorm:"index;not null;"`
	VacationPlan   VacationPlan `gorm:"foreignKey:VacationPlanID"`
	Order          int          `gorm:"not null;"`
	MemberID       uint         `gorm:"index;not null;"`
	Member         Member       `gorm:"foreignKey:MemberID"`
	DecisionDate   time.Time    `gorm:"type:date"`
}
