package models

import "cywell.com/vacation-promotion/app/models"

type NotificationType struct {
	ID            uint                  `gorm:"primaryKey"`
	TypeName      string                `gorm:"size:30"`
	Notifications []models.Notification `gorm:"foreignKey:NotificationTypeID"`
}
