package dto

import "time"

type CompanyResponse struct {
	ID                          uint      `json:"id"`
	Name                        string    `json:"name"`
	AccountingDay               time.Time `json:"accounting_day"`
	VacationGenerateTypeName    string    `json:"vacation_generate_type_name"`
	VacationGenerateDescription string    `json:"vacation_generate_description"`
}
