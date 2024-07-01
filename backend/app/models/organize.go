package models

type Organize struct {
	ID             uint        `gorm:"primaryKey" json:"organize_id"`
	CompanyID      uint        `gorm:"index" json:"company_id"`
	Name           string      `json:"organize_name"`
	ParentID       *uint       `gorm:"index" json:"parent_id"`
	ParentOrganize *Organize   `gorm:"foreignKey:ParentID" json:"parent_organize"`
	Children       []*Organize `gorm:"foreignKey:ParentID;references:ID" json:"children,omitempty"`
	Members        []*Member   `gorm:"foreignKey:OrganizeID" json:"members,omitempty"`
}
