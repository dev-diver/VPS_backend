package models

import "time"

type Member struct {
	ID                  uint    `gorm:"primaryKey"`
	CompanyID           uint    `gorm:"index"`
	Company             Company `gorm:"foreignKey:CompanyID"`
	Name                string  `gorm:"size:100"`
	Email               string  `gorm:"size:100;unique"`
	Password            string  `gorm:"size:100"`
	HireDate            time.Time
	RetireDate          *time.Time
	IsActive            bool
	Admin               []*Company            `gorm:"many2many:member_admin"`
	GivenVacations      []GivenVacation       `gorm:"foreignKey:MemberID"`
	ApplyVacations      []ApplyVacation       `gorm:"foreignKey:MemberID"`
	Groups              []*Group              `gorm:"many2many:group_members"`
	NotificationMembers []*NotificationMember `gorm:"foreignKey:MemberID"`
}
