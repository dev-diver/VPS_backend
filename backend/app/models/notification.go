package models

import "cywell.com/vacation-promotion/app/models"

type Notification struct {
	ID                  uint                    `gorm:"primaryKey"`
	NotificationTypeID  uint                    `gorm:"index"`
	NotificationType    models.NotificationType `gorm:"foreignKey:NotificationTypeID"`
	Contents            string                  `gorm:"type:text"`
	NotificationMembers []NotificationMember    `gorm:"foreignKey:NotificationID"`
}
