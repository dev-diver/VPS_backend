package models

type VacationCancelState struct {
	ID       uint   `gorm:"primaryKey"`
	TypeName string `gorm:"size:30"`
}
