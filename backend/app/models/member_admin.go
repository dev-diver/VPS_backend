package models

type MemberAdmin struct {
	CompanyID   uint      `gorm:"primaryKey"`
	MemberID    uint      `gorm:"primaryKey"`
	AdminTypeID uint      `gorm:"index"`
	AdminType   AdminType `gorm:"foreignKey:AdminTypeID"`
}
