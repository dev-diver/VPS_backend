package models

type Group struct {
	ID        uint      `gorm:"primaryKey"`
	CompanyID uint      `gorm:"index"`
	Company   Company   `gorm:"foreignKey:CompanyID"`
	Name      string    `gorm:"size:60"`
	Color     string    `gorm:"size:6"`
	Priority  int       `gorm:"default:1"`
	Members   []*Member `gorm:"many2many:group_members"`
}
