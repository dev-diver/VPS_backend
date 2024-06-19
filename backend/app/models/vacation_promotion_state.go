package models

type VacationPromotionState struct {
	ID             uint            `gorm:"primaryKey"`
	TypeName       string          `gorm:"size:30"`
	GivenVacations []GivenVacation `gorm:"foreignKey:VacationPromotionStateID"`
}
