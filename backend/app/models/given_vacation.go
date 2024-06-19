package models

import (
	"time"

	"cywell.com/vacation-promotion/app/models"
)

type GivenVacation struct {
	ID                       uint                        `gorm:"primaryKey"`
	MemberID                 uint                        `gorm:"index"`
	Member                   Member                      `gorm:"foreignKey:MemberID"`
	VacationGenerateTypeID   uint                        `gorm:"index"`
	VacationGenerateType     models.VacationGenerateType `gorm:"foreignKey:VacationGenerateTypeID"`
	VacationPromotionStateID uint                        `gorm:"index"`
	VacationPromotionState   VacationPromotionState      `gorm:"foreignKey:VacationPromotionStateID"`
	Year                     int
	GivenDays                int
	GenerateDate             time.Time
	UsedHours                int
	RemainingHours           int
	ReservedHours            int
	IsExpired                bool
}
