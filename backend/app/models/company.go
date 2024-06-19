package models

import "cywell.com/vacation-promotion/app/models"

type Company struct {
	ID                     uint                        `gorm:"primaryKey"`
	Name                   string                      `gorm:"size:60"`
	AccountingDay          string                      `gorm:"size:5"` // MM-DD 형식
	VacationGenerateTypeID uint                        `gorm:"index"`
	VacationGenerateType   models.VacationGenerateType `gorm:"foreignKey:VacationGenerateTypeID"`
	Admin                  []*Member                   `gorm:"many2many:member_admin"`
	Members                []Member                    `gorm:"foreignKey:CompanyID"`
	Groups                 []Group                     `gorm:"foreignKey:CompanyID"`
}
