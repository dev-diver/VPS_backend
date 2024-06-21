package api

import (
	"fmt"
	"strconv"
	"time"

	"cywell.com/vacation-promotion/app/dto"
	"cywell.com/vacation-promotion/app/enums"
	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/database"
	"github.com/gofiber/fiber/v2"
)

func CreateVacationPlanHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		memberID, err := strconv.ParseUint(c.Params("memberID"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid company ID"})
		}

		var request dto.CreateVacationPlanRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		vacationPlan := models.VacationPlan{
			MemberID:               uint(memberID),
			Approver1ID:            request.Approver1ID,
			ApproverFinalID:        request.ApproverFinalID,
			ApplyDate:              time.Now(),
			VacationProcessStateID: enums.VacationProcessStateApplied,
			VacationCancelStateID:  enums.VacationCancelStateDefault,
		}

		if err := db.DB.Create(&vacationPlan).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		for _, vacation := range request.Vacations {
			applyVacation := models.ApplyVacation{
				VacationPlanID:         vacationPlan.ID,
				MemberID:               uint(memberID),
				StartDate:              vacation.StartDate,
				EndDate:                vacation.EndDate,
				HalfFirst:              vacation.HalfFirst,
				HalfLast:               vacation.HalfLast,
				VacationProcessStateID: enums.VacationProcessStateApplied,
				VacationCancelStateID:  enums.VacationCancelStateDefault,
				VacationTypeID:         enums.VacationTypeNormal,
			}
			if err := db.DB.Create(&applyVacation).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		}

		//vacationPlan을 vacationPlanResponse로 변환하는 코드
		vacationPlanResponse := dto.VacationPlanResponse{
			ID:        vacationPlan.ID,
			MemberID:  vacationPlan.MemberID,
			ApplyDate: vacationPlan.ApplyDate,
		}

		for _, vacation := range request.Vacations {
			vacationPlanResponse.Vacations = append(vacationPlanResponse.Vacations, dto.ApplyVacationResponse{
				StartDate: vacation.StartDate,
				EndDate:   vacation.EndDate,
				HalfFirst: vacation.HalfFirst,
				HalfLast:  vacation.HalfLast,
			})
		}
		return c.Status(fiber.StatusCreated).JSON(vacationPlanResponse)
	}
}

func GetVacationsByPeriodHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		companyID, groupID, memberID, year, month, err := parseParams(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		var startDate, endDate time.Time
		if month != 0 {
			startDate = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			endDate = startDate.AddDate(0, 1, -1)
		} else {
			startDate = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
			endDate = startDate.AddDate(1, 0, -1)
		}

		var vacations []models.ApplyVacation
		query := db.DB.Preload("Member")

		if companyID != 0 {
			query = query.Joins("JOIN members ON members.id = apply_vacations.member_id").Where("members.company_id = ?", companyID)
		} else if groupID != 0 {
			query = query.Joins("JOIN group_members ON group_members.member_id = apply_vacations.member_id").
				Joins("JOIN members ON members.id = group_members.member_id").
				Where("group_members.group_id = ?", groupID)
		} else if memberID != 0 {
			query = query.Where("member_id = ?", memberID)
		}

		query = query.Where("start_date <= ? AND end_date >= ?", endDate, startDate)

		if err := query.Find(&vacations).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		vacationsResponse := make([]dto.ApplyVacationCardResponse, 0)
		for _, vacation := range vacations {
			vacationsResponse = append(vacationsResponse, dto.ApplyVacationCardResponse{
				ID:           vacation.ID,
				MemberID:     vacation.MemberID,
				MemberName:   vacation.Member.Name,
				StartDate:    vacation.StartDate,
				EndDate:      vacation.EndDate,
				HalfFirst:    vacation.HalfFirst,
				HalfLast:     vacation.HalfLast,
				Status:       vacation.VacationProcessStateID,
				CancelStatus: vacation.VacationCancelStateID,
			})
		}
		return c.JSON(vacationsResponse)
	}
}

func GetVacationPlansByPeriodHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		companyID, groupID, memberID, year, month, err := parseParams(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		var startDate, endDate time.Time
		if month != 0 {
			startDate = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			endDate = startDate.AddDate(0, 1, -1)
		} else {
			startDate = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
			endDate = startDate.AddDate(1, 0, -1)
		}

		var vacationPlans []models.VacationPlan
		query := db.DB.Preload("ApplyVacations.Member").
			Preload("Member").
			Joins("JOIN apply_vacations ON apply_vacations.vacation_plan_id = vacation_plans.id").
			Joins("JOIN members ON members.id = vacation_plans.member_id").
			Where("apply_vacations.start_date <= ? AND apply_vacations.end_date >= ?", endDate, startDate)

		if companyID != 0 {
			query = query.Where("members.company_id = ?", companyID)
		} else if groupID != 0 {
			query = query.Joins("JOIN group_members ON group_members.member_id = vacation_plans.member_id").
				Where("group_members.group_id = ?", groupID)
		} else if memberID != 0 {
			query = query.Where("vacation_plans.member_id = ?", memberID)
		}

		if err := query.Find(&vacationPlans).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		response := make([]dto.VacationPlanResponse, 0)
		for _, plan := range vacationPlans {
			var vacations []dto.ApplyVacationResponse
			var earliestDate, latestDate time.Time
			var withinRange bool
			for i, vacation := range plan.ApplyVacations {
				if vacation.StartDate.Before(endDate) && vacation.EndDate.After(startDate) {
					withinRange = true
					vacations = append(vacations, dto.ApplyVacationResponse{
						ID:           vacation.ID,
						StartDate:    vacation.StartDate,
						EndDate:      vacation.EndDate,
						HalfFirst:    vacation.HalfFirst,
						HalfLast:     vacation.HalfLast,
						Status:       vacation.VacationProcessStateID,
						CancelStatus: vacation.VacationCancelStateID,
					})
					if i == 0 || vacation.StartDate.Before(earliestDate) {
						earliestDate = vacation.StartDate
					}
					if i == 0 || vacation.EndDate.After(latestDate) {
						latestDate = vacation.EndDate
					}
				}
			}
			if withinRange {
				response = append(response, dto.VacationPlanResponse{
					ID:         plan.ID,
					MemberID:   plan.MemberID,
					MemberName: plan.Member.Name,
					ApplyDate:  plan.ApplyDate,
					Vacations:  vacations,
				})
			}
		}

		return c.JSON(response)
	}
}

func ApproveVacationPlanHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

func EditVacationPlanHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

func UpdateVacationHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

func CancelVacationHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

func ApproveVacationHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

func RejectVacationHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

func GetPromotionsHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		companyID := c.Params("companyID")
		var promotions []models.GivenVacation
		if err := db.DB.Where("company_id = ? AND vacation_promotion_state_id = ?", companyID, 2).Find(&promotions).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(promotions)
	}
}

func parseParams(c *fiber.Ctx) (uint64, uint64, uint64, int, int, error) {
	companyIDStr := c.Params("companyID")
	groupIDStr := c.Params("groupID")
	memberIDStr := c.Params("memberID")
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	var companyID, groupID, memberID uint64
	var year, month int
	var err error

	if companyIDStr != "" {
		companyID, err = strconv.ParseUint(companyIDStr, 10, 32)
		if err != nil {
			return 0, 0, 0, 0, 0, fmt.Errorf("Invalid company ID")
		}
	}
	if groupIDStr != "" {
		groupID, err = strconv.ParseUint(groupIDStr, 10, 32)
		if err != nil {
			return 0, 0, 0, 0, 0, fmt.Errorf("Invalid group ID")
		}
	}
	if memberIDStr != "" {
		memberID, err = strconv.ParseUint(memberIDStr, 10, 32)
		if err != nil {
			return 0, 0, 0, 0, 0, fmt.Errorf("Invalid member ID")
		}
	}

	year, err = strconv.Atoi(yearStr)
	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("Invalid year")
	}

	if monthStr != "" {
		month, err = strconv.Atoi(monthStr)
		if err != nil {
			return 0, 0, 0, 0, 0, fmt.Errorf("Invalid month")
		}
	}

	return companyID, groupID, memberID, year, month, nil
}
