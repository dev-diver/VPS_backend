package models

import "time"

type Member struct {
	ID               uint    `gorm:"primaryKey"`
	CompanyID        uint    `gorm:"index"`
	Company          Company `gorm:"foreignKey:CompanyID"`
	Name             string  `gorm:"size:100"`
	Email            string  `gorm:"size:100;unique"`
	Password         string  `gorm:"size:100"`
	HireDate         time.Time
	RetireDate       *time.Time
	IsActive         bool
	Admin            []*Company           `gorm:many2many:member_admin`
	GivenVacations   []GivenVacation      `gorm:"foreignKey:MemberID"`
	ConsumeVacations []ConsumeVacation    `gorm:"foreignKey:MemberID"`
	GroupMembers     []GroupMember        `gorm:"foreignKey:MemberID"`
	Notifications    []NotificationMember `gorm:"foreignKey:MemberID"`
}
