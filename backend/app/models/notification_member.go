package models

type NotificationMember struct {
	MemberID       uint         `gorm:"primaryKey"`
	Member         Member       `gorm:"foreignKey:MemberID"`
	NotificationID uint         `gorm:"primaryKey"`
	Notification   Notification `gorm:"foreignKey:NotificationID"`
	IsApprove      bool
}
