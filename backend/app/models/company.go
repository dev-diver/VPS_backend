package models

import "time"

type Company struct {
	ID                     uint                 `gorm:"primaryKey"`
	Name                   string               `gorm:"size:60"`
	AccountingDay          time.Time            `gorm:"type:date"` // MM-DD 형식
	VacationGenerateTypeID uint                 `gorm:"index"`
	VacationGenerateType   VacationGenerateType `gorm:"foreignKey:VacationGenerateTypeID"`
	Admin                  []*Member            `gorm:"many2many:member_admin"`
	Members                []Member             `gorm:"foreignKey:CompanyID"`
	Groups                 []Group              `gorm:"foreignKey:CompanyID"`
}
