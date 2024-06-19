package models

type Notification struct {
	ID                  uint                  `gorm:"primaryKey"`
	NotificationTypeID  uint                  `gorm:"index"`
	NotificationType    NotificationType      `gorm:"foreignKey:NotificationTypeID"`
	Contents            string                `gorm:"type:text"`
	NotificationMembers []*NotificationMember `gorm:"foreignKey:NotificationID"`
}
