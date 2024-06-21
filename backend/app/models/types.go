package models

type VacationType struct {
	ID       uint   `gorm:"primaryKey"`
	TypeName string `gorm:"size:30"`
}

type VacationPromotionState struct {
	ID       uint   `gorm:"primaryKey"`
	TypeName string `gorm:"size:30"`
}

type VacationProcessState struct {
	ID       uint   `gorm:"primaryKey"`
	TypeName string `gorm:"size:30"`
}

type VacationGenerateType struct {
	ID          uint   `gorm:"primaryKey"`
	TypeName    string `gorm:"size:30"`
	Description string `gorm:"type:text"`
}

type VacationCancelState struct {
	ID       uint   `gorm:"primaryKey"`
	TypeName string `gorm:"size:30"`
}

type NotificationType struct {
	ID       uint   `gorm:"primaryKey"`
	TypeName string `gorm:"size:30"`
}

type AdminType struct {
	ID       uint   `gorm:"primaryKey"`
	TypeName string `gorm:"size:30"`
}
