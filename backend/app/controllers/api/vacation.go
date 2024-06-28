package api

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"cywell.com/vacation-promotion/app/dto"
	"cywell.com/vacation-promotion/app/enums"
	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/database"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

		// 데이터 유효성 검사
		validate := validator.New()
		if err := validate.Struct(request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		vacationPlan := models.VacationPlan{
			MemberID:     uint(memberID),
			ApplyDate:    time.Now(),
			ApproveStage: 0,
			RejectState:  false,
		}

		err = db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&vacationPlan).Error; err != nil {
				return err
			}

			for i, approverID := range request.Approvers {
				approverOrder := models.ApproverOrder{
					VacationPlanID: vacationPlan.ID,
					Order:          i + 1,
					MemberID:       uint(approverID),
				}
				if err := tx.Create(&approverOrder).Error; err != nil {
					return err
				}
			}

			for _, vacation := range request.Vacations {
				applyVacation := models.ApplyVacation{
					VacationPlanID: vacationPlan.ID,
					MemberID:       uint(memberID),
					StartDate:      vacation.StartDate,
					EndDate:        vacation.EndDate,
					HalfFirst:      vacation.HalfFirst,
					HalfLast:       vacation.HalfLast,
					ApproveStage:   0,
					RejectState:    false,
					VacationTypeID: enums.VacationTypeNormal,
				}
				if err := tx.Create(&applyVacation).Error; err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

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

func GetVacationPlanHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		planID, err := strconv.ParseUint(c.Params("planID"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid plan ID"})
		}

		var plan models.VacationPlan //TODO : Preload 최적화
		if err := db.DB.
			Preload("Member").
			Preload("ApplyVacations").
			Preload("ApproverOrders").
			Preload("ApproverOrders.Member").
			First(&plan, planID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Vacation plan not found"})
		}

		vacationPlanResponse := dto.MapVacationPlanToResponse(plan)

		for _, ApproveOrder := range plan.ApproverOrders {
			vacationPlanResponse.ApproverOrder = append(vacationPlanResponse.ApproverOrder, dto.MapApproverOrderToResponse(ApproveOrder))
		}

		for _, vacation := range plan.ApplyVacations {
			vacationPlanResponse.Vacations = append(vacationPlanResponse.Vacations, dto.MapApplyVacationToResponse(vacation))
		}

		return c.JSON(vacationPlanResponse)
	}
}

func GetVacationHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		vacationID, err := strconv.ParseUint(c.Params("vacationID"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid plan ID"})
		}

		var vacation models.ApplyVacation
		if err := db.DB.Preload("Member").First(&vacation, vacationID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Vacation not found"})
		}

		vacationResponse := dto.MapApplyVacationToResponse(vacation)
		return c.JSON(vacationResponse)
	}
}

func GetVacationsByPeriodHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		companyID, groupID, memberID, _, year, month, err := parseParams(c)
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
			vacationsResponse = append(vacationsResponse, dto.MapApplyVacationToCardResponse(vacation))
		}
		return c.JSON(vacationsResponse)
	}
}

func GetVacationPlansByPeriodHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		companyID, groupID, memberID, approverID, year, month, err := parseParams(c)
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

		var vacationPlans []models.VacationPlan //TODO : Preload 최적화
		query := db.DB.Preload("ApplyVacations", "start_date <= ? AND end_date >= ?", endDate, startDate)

		if companyID != 0 {
			query = query.Joins("JOIN members ON members.id = vacation_plans.member_id").
				Where("members.company_id = ?", companyID).
				Preload("Member", "company_id = ?", companyID)
		} else if groupID != 0 {
			query = query.Joins("JOIN group_members ON group_members.member_id = vacation_plans.member_id").
				Joins("JOIN members ON members.id = group_members.member_id").
				Where("group_members.group_id = ?", groupID).
				Preload("Member", "id IN (SELECT member_id FROM group_members WHERE group_id = ?)", groupID)
		} else if memberID != 0 {
			query = query.Where("vacation_plans.member_id = ?", memberID).
				Preload("Member")
		} else if approverID != 0 {
			query = query.Joins("JOIN approver_orders ON approver_orders.vacation_plan_id = vacation_plans.id").
				Where("approver_orders.member_id = ?", approverID).
				Preload("Member")
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid query"})
		}

		if err := query.
			Preload("ApproverOrders").
			Preload("ApproverOrders.Member").
			Find(&vacationPlans).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		response := make([]dto.VacationPlanResponse, 0)
		for _, plan := range vacationPlans {

			vacationPlanResponse := dto.MapVacationPlanToResponse(plan)

			for _, vacation := range plan.ApplyVacations {
				vacationPlanResponse.Vacations = append(vacationPlanResponse.Vacations, dto.MapApplyVacationToResponse(vacation))
			}

			for _, ApproveOrder := range plan.ApproverOrders {
				vacationPlanResponse.ApproverOrder = append(vacationPlanResponse.ApproverOrder, dto.MapApproverOrderToResponse(ApproveOrder))
			}

			response = append(response, vacationPlanResponse)

		}

		return c.JSON(response)
	}
}

func ApproveVacationPlanHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {

		//ID 이상 검증
		planID, err := strconv.ParseUint(c.Params("planID"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "잘못된 계획 ID입니다"})
		}

		//Vacation 검증
		input, plan, err := ValidateVacationPlanRequest(c, db, planID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := ValdiateVacationPlanApproval(c, db, input, plan); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		//거절되었는지 확인
		if plan.RejectState {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "휴가 계획이 거절 상태입니다."})
		}

		// 휴가 계획 상태 업데이트
		plan.ApproveStage = uint(input.ApprovalStage)
		if err := db.DB.Save(&plan).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "휴가 계획을 승인할 수 없습니다"})
		}

		// 휴가 상태 업데이트
		for _, vacation := range plan.ApplyVacations {
			if !vacation.RejectState {
				vacation.ApproveStage = uint(input.ApprovalStage)
				if err := db.DB.Save(&vacation).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "휴가 상태를 업데이트할 수 없습니다"})
				}
			}
		}

		return c.JSON(plan)
	}
}

func RejectVacationPlanHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {

		//ID 이상
		planID, err := strconv.ParseUint(c.Params("planID"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "잘못된 계획 ID입니다"})
		}

		//Vacation 검증
		input, plan, err := ValidateVacationPlanRequest(c, db, planID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := ValdiateVacationPlanApproval(c, db, input, plan); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		//거절되었는지 확인
		if plan.RejectState {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "휴가 계획이 거절되었습니다"})
		}

		// 휴가 계획 상태 업데이트
		plan.RejectState = true
		if err := db.DB.Save(&plan).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "휴가 계획을 거절할 수 없습니다"})
		}

		// 휴가 상태 업데이트
		for _, vacation := range plan.ApplyVacations {
			plan.RejectState = true
			if err := db.DB.Save(&vacation).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "휴가 상태를 업데이트할 수 없습니다"})
			}
		}

		return c.JSON(plan)
	}
}

func CancelRejectVacationPlanHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {

		//ID 이상
		planID, err := strconv.ParseUint(c.Params("planID"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "잘못된 계획 ID입니다"})
		}

		//Vacation 검증
		input, plan, err := ValidateVacationPlanRequest(c, db, planID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := ValdiateVacationPlanApproval(c, db, input, plan); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// 휴가 계획 상태 업데이트
		plan.RejectState = false
		if err := db.DB.Save(&plan).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "휴가 계획을 거절할 수 없습니다"})
		}

		// 휴가 상태 업데이트
		for _, vacation := range plan.ApplyVacations {
			plan.RejectState = false
			if err := db.DB.Save(&vacation).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "휴가 상태를 업데이트할 수 없습니다"})
			}
		}

		return c.JSON(plan)
	}
}

func UpdateVacationPlanHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

func UpdateVacationHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// 요청자
func DeleteVacationPlanHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		planID := c.Params("planID")
		var plan models.VacationPlan

		tx := db.DB.Begin()
		if err := tx.Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if err := tx.Preload("ApplyVacations").First(&plan, planID).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if err := tx.Where("vacation_plan_id = ?", planID).Delete(&models.ApplyVacation{}).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if err := tx.Delete(&plan).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Vacation plan deleted successfully"})
	}
}

func DeleteVacationHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		vacationId := c.Params("vacationID")
		var vacation models.ApplyVacation
		if err := db.DB.First(&vacation, vacationId).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if err := db.DB.Delete(&vacation).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Vacation deleted successfully"})
	}
}

// 결재자
func RejectVacationHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		//TODO:결재자만 reject 가능
		vacationId := c.Params("vacationID")
		var vacation models.ApplyVacation
		if err := db.DB.First(&vacation, vacationId).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		vacation.RejectState = true
		if err := db.DB.Save(&vacation).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		vacationResponse := dto.MapApplyVacationToResponse(vacation)
		return c.JSON(vacationResponse)
	}
}

func CancelRejectVacationHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		//TODO:결재자만 원복 가능
		vacationId := c.Params("vacationID")
		var vacation models.ApplyVacation
		if err := db.DB.First(&vacation, vacationId).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		vacation.RejectState = false
		if err := db.DB.Save(&vacation).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		vacationResponse := dto.MapApplyVacationToResponse(vacation)
		return c.JSON(vacationResponse)
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

func parseParams(c *fiber.Ctx) (uint64, uint64, uint64, uint64, int, int, error) {
	companyIDStr := c.Params("companyID")
	groupIDStr := c.Params("groupID")
	memberIDStr := c.Params("memberID")
	yearStr := c.Query("year")
	monthStr := c.Query("month")
	approverIDstr := c.Query("approverID")

	var companyID, groupID, memberID, approverID uint64
	var year, month int
	var err error

	if companyIDStr != "" {
		companyID, err = strconv.ParseUint(companyIDStr, 10, 32)
		if err != nil {
			return 0, 0, 0, 0, 0, 0, fmt.Errorf("invalid company ID")
		}
	}
	if groupIDStr != "" {
		groupID, err = strconv.ParseUint(groupIDStr, 10, 32)
		if err != nil {
			return 0, 0, 0, 0, 0, 0, fmt.Errorf("invalid group ID")
		}
	}
	if memberIDStr != "" {
		memberID, err = strconv.ParseUint(memberIDStr, 10, 32)
		if err != nil {
			return 0, 0, 0, 0, 0, 0, fmt.Errorf("invalid member ID")
		}
	}
	if approverIDstr != "" {
		approverID, err = strconv.ParseUint(approverIDstr, 10, 32)
		if err != nil {
			return 0, 0, 0, 0, 0, 0, fmt.Errorf("invalid approver ID")
		}
	}

	year, err = strconv.Atoi(yearStr)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, fmt.Errorf("invalid year")
	}

	if monthStr != "" {
		month, err = strconv.Atoi(monthStr)
		if err != nil {
			return 0, 0, 0, 0, 0, 0, fmt.Errorf("invalid month")
		}
	}

	return companyID, groupID, memberID, approverID, year, month, nil
}

func ValidateVacationPlanRequest(c *fiber.Ctx, db *database.Database, planID uint64) (dto.ApproveVacationPlanRequest, models.VacationPlan, error) {

	// 요청 바디 검증
	input := dto.ApproveVacationPlanRequest{}
	plan := models.VacationPlan{}
	if err := c.BodyParser(&input); err != nil {
		return input, plan, err
	}

	// 요청 데이터 검증
	validate := validator.New()
	if err := validate.Struct(&input); err != nil {
		return input, plan, err
	}

	if err := db.Preload("ApplyVacations").First(&plan, planID).Error; err != nil {
		return input, plan, err
	}

	return input, plan, nil
}

func ValdiateVacationPlanApproval(c *fiber.Ctx, db *database.Database, input dto.ApproveVacationPlanRequest, plan models.VacationPlan) error {

	// 승인 단계가 올바른지 검증
	if input.ApprovalStage <= plan.ApproveStage && !plan.RejectState {
		return errors.New("잘못된 승인 단계 순서입니다")
	}

	// 지정된 승인자가 승인하는지 검증
	var expectedMemberID uint
	if err := db.Table("approver_orders").Where("vacation_plan_id = ? AND `order` = ?", plan.ID, input.ApprovalStage).Pluck("member_id", &expectedMemberID).Error; err != nil {
		return errors.New("승인 권한을 찾을 수 없습니다")
	}

	if expectedMemberID != input.MemberID {
		return errors.New("승인 권한이 없습니다")
	}

	return nil
}
