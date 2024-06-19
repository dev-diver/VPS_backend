package models

type VacationProcessState struct {
	ID       uint   `gorm:"primaryKey"`
	TypeName string `gorm:"size:30"`
}
