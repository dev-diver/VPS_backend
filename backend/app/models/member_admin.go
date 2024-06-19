package models

import "cywell.com/vacation-promotion/app/models"

type MemberAdmin struct {
	CompanyID   uint             `gorm:"primaryKey"`
	MemberID    uint             `gorm:"primaryKey"`
	AdminTypeID uint             `gorm:"index"`
	AdminType   models.AdminType `gorm:"foreignKey:AdminTypeID"`
}
