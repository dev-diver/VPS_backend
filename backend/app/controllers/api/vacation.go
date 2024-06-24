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

func GetVacationPlanHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		planID, err := strconv.ParseUint(c.Params("planID"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid plan ID"})
		}

		var plan models.VacationPlan
		if err := db.DB.Preload("Member").Preload("ApplyVacations").First(&plan, planID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Vacation plan not found"})
		}

		vacationPlanResponse := dto.MapVacationPlanToResponse(plan)

		for _, vacation := range plan.ApplyVacations {
			vacationPlanResponse.Vacations = append(vacationPlanResponse.Vacations, dto.MapApplyVacationToResponse(vacation))
		}

		return c.JSON(vacationPlanResponse)
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
			vacationsResponse = append(vacationsResponse, dto.MapApplyVacationToCardResponse(vacation))
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
		query := db.DB.Preload("ApplyVacations", "start_date <= ? AND end_date >= ?", endDate, startDate).
			Preload("ApplyVacations.Member")

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
		}

		if err := query.Find(&vacationPlans).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		response := make([]dto.VacationPlanResponse, 0)
		for _, plan := range vacationPlans {

			vacationPlanResponse := dto.MapVacationPlanToResponse(plan)

			for _, vacation := range plan.ApplyVacations {
				vacationPlanResponse.Vacations = append(vacationPlanResponse.Vacations, dto.MapApplyVacationToResponse(vacation))
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

		if err := ValdiateVacationPlanApproval(c, input, plan); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		//거절되었는지 확인
		if plan.VacationProcessStateID == enums.VacationProcessStateRejected {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "휴가 계획이 거절 상태입니다."})
		}

		// 휴가 계획 상태 업데이트
		plan.VacationProcessStateID = uint(input.ApprovalState)
		if err := db.DB.Save(&plan).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "휴가 계획을 승인할 수 없습니다"})
		}

		// 휴가 상태 업데이트
		for _, vacation := range plan.ApplyVacations {
			if vacation.VacationProcessStateID != enums.VacationProcessStateRejected {
				vacation.VacationProcessStateID = uint(input.ApprovalState)
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

		if err := ValdiateVacationPlanApproval(c, input, plan); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		//거절되었는지 확인
		if plan.VacationProcessStateID == enums.VacationProcessStateRejected {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "휴가 계획이 거절되었습니다"})
		}

		// 휴가 계획 상태 업데이트
		plan.VacationProcessStateID = enums.VacationProcessStateRejected
		if err := db.DB.Save(&plan).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "휴가 계획을 거절할 수 없습니다"})
		}

		// 휴가 상태 업데이트
		for _, vacation := range plan.ApplyVacations {
			vacation.VacationProcessStateID = enums.VacationProcessStateRejected
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

		if err := ValdiateVacationPlanApproval(c, input, plan); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// 휴가 계획 상태 업데이트
		plan.VacationProcessStateID = input.ApprovalState - 1
		if err := db.DB.Save(&plan).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "휴가 계획을 거절할 수 없습니다"})
		}

		// 휴가 상태 업데이트
		for _, vacation := range plan.ApplyVacations {
			vacation.VacationProcessStateID = input.ApprovalState - 1
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
func DeleteVacationHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
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
		vacation.VacationProcessStateID = enums.VacationProcessStateRejected
		vacation.VacationCancelStateID = enums.VacationCancelStateCompleted
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
		vacation.VacationProcessStateID = enums.VacationProcessStateApplied //TODO: 이전 레벨로
		vacation.VacationCancelStateID = enums.VacationCancelStateDefault
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

	if err := db.DB.Preload("ApplyVacations").First(&plan, planID).Error; err != nil {
		return input, plan, err
	}

	return input, plan, nil
}

func ValdiateVacationPlanApproval(c *fiber.Ctx, input dto.ApproveVacationPlanRequest, plan models.VacationPlan) error {
	// 계획 ID 검증

	// 승인 단계가 올바른지 검증
	if input.ApprovalState <= plan.VacationProcessStateID && plan.VacationProcessStateID != enums.VacationProcessStateRejected {
		return errors.New("잘못된 승인 단계 순서입니다")
	}

	// 지정된 승인자가 승인하는지 검증
	var expectedMemberID uint
	if input.ApprovalState == enums.VacationProcessStateFirstApproved {
		expectedMemberID = plan.Approver1ID
	} else if input.ApprovalState == enums.VacationProcessStateFinalApproved {
		expectedMemberID = plan.ApproverFinalID
	} else {
		return errors.New("잘못된 승인자입니다")
	}

	if expectedMemberID != input.MemberID {
		return errors.New("승인 권한이 없습니다")
	}
	// 여기까지 통과했다면 검증 성공
	return nil
}
