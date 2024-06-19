package models

import "cywell.com/vacation-promotion/app/models"

type VacationType struct {
	ID               uint                     `gorm:"primaryKey"`
	TypeName         string                   `gorm:"size:30"`
	ConsumeVacations []models.ConsumeVacation `gorm:"foreignKey:VacationTypeID"`
}
