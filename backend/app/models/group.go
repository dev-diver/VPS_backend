package models

type Group struct {
	ID           uint    `gorm:"primaryKey"`
	CompanyID    uint    `gorm:"index"`
	Company      Company `gorm:"foreignKey:CompanyID"`
	Name         string  `gorm:"size:60"`
	Color        string  `gorm:"size:6"`
	Priority     int
	GroupMembers []GroupMember `gorm:"foreignKey:GroupID"`
}
