package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"cywell.com/vacation-promotion/app/enums"
	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/app/utils"
	"cywell.com/vacation-promotion/database"
	"cywell.com/vacation-promotion/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	app := fiber.New()

	hostIP, err := getLocalIP()
	if err != nil {
		log.Fatalf("Failed to get local IP: %v", err)
	}

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost") ||
				strings.HasPrefix(origin, "http://"+hostIP)
		},
		AllowCredentials: true,
	}))

	if err := utils.SetJWTSecretKey(); err != nil {
		fmt.Println(err.Error())
		log.Fatal("failed to set jwt secret")
	}

	db, err := database.InitDB()
	if err != nil {
		log.Fatal("failed to connect database")
	}
	fmt.Println("Database seeding successfully")

	err = MigrateAndSeed(db)
	if err != nil {
		log.Fatal("failed to migrate types: ", err)
	}
	err = db.AutoMigrate(
		&models.Company{},
		&models.Member{},
		&models.MemberAdmin{},
		&models.NotificationMember{},
		&models.Group{},
		&models.GivenVacation{},
		&models.ApplyVacation{},
		&models.VacationPlan{},
		&models.Notification{},
		&models.ApproverOrder{},
		&models.Organize{},
	)

	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}
	fmt.Println("Database migrated successfully")

	api := app.Group("/api")
	routes.RegisterAPI(api, db)

	app.Static("/", "../dist/front_web/browser/")

	app.Use(func(c *fiber.Ctx) error {
		// Return index.html for all other routes
		if err := c.SendFile("../dist/front_web/browser/index.html"); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return nil
	})

	app.Listen(":3000")
}

func MigrateAndSeed(db *database.Database) error {
	err := db.AutoMigrate(
		&models.VacationType{},
		&models.VacationPromotionState{},
		&models.VacationGenerateType{},
		&models.NotificationType{},
		&models.AdminType{},
	)

	if err != nil {
		return err
	}

	//관리자 타입
	adminTypes := []models.AdminType{
		{ID: enums.AdminTypeManager, TypeName: "관리자"},
	}
	for _, at := range adminTypes {
		db.FirstOrCreate(&at, models.AdminType{ID: at.ID})
	}

	//휴가 타입
	vacationTypes := []models.VacationType{
		{ID: enums.VacationTypeNormal, TypeName: "연차"},
		{ID: enums.VacationTypePromotion, TypeName: "촉진"},
	}
	for _, vt := range vacationTypes {
		db.FirstOrCreate(&vt, models.VacationType{ID: vt.ID})
	}

	//휴가 발생 타입
	vacationGenerateTypes := []models.VacationGenerateType{
		{ID: enums.VacationGenerateTypeAnnualNormal, TypeName: "입사일 지급-월1일", Description: "1년차까지 매월 1일 지급, 이후 년차별로 지급"},
		{ID: enums.VacationGenerateTypeAnnualThisYearPreGiven, TypeName: "입사일 지급-당해년도 선지급", Description: "입사한 해 남은 달만큼 선지급, 이후 년차별로 지급"},
		{ID: enums.VacationGenerateTypeAnnualOneYearPreGiven, TypeName: "입사일지급-11일 선지급", Description: "11일 선지급, 2년차부터 년차별로 지급"},

		{ID: enums.VacationGenerateTypePreAccountingNormal, TypeName: "회계일 지급-월1일-선지급", Description: "입사한 해 매월 1일 지급, 다음해 회계일기준 선지급"},
		{ID: enums.VacationGenerateTypePreAccountingThisYearPreGiven, TypeName: "회계일 지급-당해년도 선지급-선지급", Description: "입사한 해 남은 달만큼 선지급, 다음해 회계일기준 선지급"},
		{ID: enums.VacationGenerateTypePreAccountingOneYearPreGiven, TypeName: "회계일 지급-11일 선지급-선지급", Description: "11일 선지급, 입사 다음해 회계일기준 선지급"},

		{ID: enums.VacationGenerateTypeProAccountingNormal, TypeName: "회계일 지급-월1일-비례지급", Description: "입사한 해 매월 1일 지급, 다음해 회계일기준 비례지급"},
		{ID: enums.VacationGenerateTypeProAccountingThisYearPreGiven, TypeName: "회계일 지급-당해년도 선지급-비례지급", Description: "입사한 해 남은 달만큼 선지급, 다음해 회계일기준 비례지급"},
		{ID: enums.VacationGenerateTypeProAccountingOneYearPreGiven, TypeName: "회계일 지급-11일 선지급-비례지급", Description: "11일 선지급, 입사 다음해 회계일기준 비례지급"},
	}
	for _, vgt := range vacationGenerateTypes {
		db.FirstOrCreate(&vgt, models.VacationGenerateType{ID: vgt.ID})
	}

	//휴가 촉진 상태
	vacationPromotionStates := []models.VacationPromotionState{
		{ID: enums.VacationPromotionStateNone, TypeName: "촉진없음"},
		{ID: enums.VacationPromotionStateFirstNoti, TypeName: "1차촉진전송"},
		{ID: enums.VacationPromotionStateFirstComplete, TypeName: "1차촉진완료"},
		{ID: enums.VacationPromotionStateSecondNeed, TypeName: "2차촉진필요"},
		{ID: enums.VacationPromotionStateSecondNoti, TypeName: "2차촉진전송"},
		{ID: enums.VacationPromotionStateSecondComplete, TypeName: "2차촉진완료"},
	}
	for _, vps := range vacationPromotionStates {
		db.FirstOrCreate(&vps, models.VacationPromotionState{ID: vps.ID})
	}

	//알림 타입
	notificationTypes := []models.NotificationType{
		{ID: enums.NotificationTypeNormal, TypeName: "일반"},
		{ID: enums.NotificationTypeVacationApplied, TypeName: "휴가 신청"},
		{ID: enums.NotificationTypeVacationFirstPromotion, TypeName: "1차 촉진"},
		{ID: enums.NotificationTypeVacationFirstPromotionAccept, TypeName: "1차 촉진 확인"},
		{ID: enums.NotificationTypeVacationSecondPromotion, TypeName: "2차 촉진"},
		{ID: enums.NotificationTypeVacationSecondPromotionAccept, TypeName: "2차 촉진 확인"},
		{ID: enums.NotificationTypeVacationDenyWork, TypeName: "노무 거부"},
		{ID: enums.NotificationTypeVacationDenyWorkAccept, TypeName: "노무 거부 확인"},
	}
	for _, nt := range notificationTypes {
		db.FirstOrCreate(&nt, models.NotificationType{ID: nt.ID})
	}

	return nil
}

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", nil
}
