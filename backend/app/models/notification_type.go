package models

type NotificationType struct {
	ID       uint   `gorm:"primaryKey"`
	TypeName string `gorm:"size:30"`
}
