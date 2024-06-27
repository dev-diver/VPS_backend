package models

import (
	"time"
)

type ApproverOrder struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	VacationPlanID uint      `gorm:"not null"`
	Order          int       `gorm:"not null;"`
	MemberID       uint      `gorm:"not null;"`
	Member         Member    `gorm:"foreignKey:MemberID"`
	DecisionDate   time.Time `gorm:"type:date"`
}
