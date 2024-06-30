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
	Admin               []*Company            `gorm:"many2many:member_admins"`
	GivenVacations      []GivenVacation       `gorm:"foreignKey:MemberID"`
	ApplyVacations      []ApplyVacation       `gorm:"foreignKey:MemberID"`
	VacationPlans       []VacationPlan        `gorm:"foreignKey:MemberID"`
	ApproverOrders      []ApproverOrder       `gorm:"foreignKey:MemberID"`
	OrganizeID          *uint                 `gorm:"index"`
	Organize            *Organize             `gorm:"foreignKey:OrganizeID"`
	Groups              []*Group              `gorm:"many2many:group_members"`
	NotificationMembers []*NotificationMember `gorm:"foreignKey:MemberID"`
}
