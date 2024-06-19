package models

type GroupMember struct {
	GroupID  uint   `gorm:"primaryKey"`
	Group    Group  `gorm:"foreignKey:GroupID"`
	MemberID uint   `gorm:"primaryKey"`
	Member   Member `gorm:"foreignKey:MemberID"`
}
