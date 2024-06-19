package models

type VacationGenerateType struct {
	ID          uint   `gorm:"primaryKey"`
	TypeName    string `gorm:"size:30"`
	Description string `gorm:"type:text"`
}
