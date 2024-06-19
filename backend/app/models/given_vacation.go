package models

import (
	"time"
)

type GivenVacation struct {
	ID                       uint                   `gorm:"primaryKey"`
	MemberID                 uint                   `gorm:"index"`
	Member                   Member                 `gorm:"foreignKey:MemberID"`
	VacationGenerateTypeID   uint                   `gorm:"index"`
	VacationGenerateType     VacationGenerateType   `gorm:"foreignKey:VacationGenerateTypeID"`
	VacationPromotionStateID uint                   `gorm:"index"`
	VacationPromotionState   VacationPromotionState `gorm:"foreignKey:VacationPromotionStateID"`
	Year                     int
	GivenDays                float32
	GenerateDate             time.Time
	UsedDays                 int
	RemainingDays            int
	ReservedDays             int
	IsExpired                bool
}
