package models

type VacationGenerateType struct {
	ID             uint            `gorm:"primaryKey"`
	TypeName       string          `gorm:"size:30"`
	Description    string          `gorm:"type:text"`
	Companies      []Company       `gorm:"foreignKey:VacationGenerateTypeID"`
	GivenVacations []GivenVacation `gorm:"foreignKey:VacationGenerateTypeID"`
}
