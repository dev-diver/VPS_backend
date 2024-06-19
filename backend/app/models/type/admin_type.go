package models

import "cywell.com/vacation-promotion/app/models"

type AdminType struct {
	ID           uint                 `gorm:"primaryKey"`
	TypeName     string               `gorm:"size:30"`
	MemberAdmins []models.MemberAdmin `gorm:"foreignKey:AdminTypeID"`
}
