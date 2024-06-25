package routes

import (
	"cywell.com/vacation-promotion/app/controllers/api"
	"cywell.com/vacation-promotion/database"
	"github.com/gofiber/fiber/v2"
)

func RegisterAPI(apiRouter fiber.Router, db *database.Database) {
	registerAuth(apiRouter, db)
	registerCompanies(apiRouter, db)
	registerGroups(apiRouter, db)
	registerMembers(apiRouter, db)
	registerVacations(apiRouter, db)
}

func registerAuth(apiRouter fiber.Router, db *database.Database) {
	auth := apiRouter.Group("/auth")
	auth.Post("/login", api.LoginHandler(db))
}

func registerCompanies(apiRouter fiber.Router, db *database.Database) {
	companies := apiRouter.Group("/companies")
	companies.Post("/", api.CreateCompanyHandler(db))

	company := companies.Group("/:companyID")
	company.Get("/", api.GetCompanyHandler(db))
	company.Post("/", api.UpdateCompanyHandler(db))
	company.Delete("/", api.DeleteCompanyHandler(db))

	members := company.Group("/members")
	members.Get("/", api.GetCompanyMembersHandler(db))
	members.Get("/search", api.SearchMembersHandler(db))   // keyword
	members.Post("/", api.CreateCompanyMembersHandler(db)) // []

	groups := company.Group("/groups")
	groups.Get("/", api.GetGroupsHandler(db))
	groups.Post("/", api.CreateGroupHandler(db))

	vacations := company.Group("/vacations")
	vacations.Get("/", api.GetVacationsByPeriodHandler(db))
	vacations.Get("/plans", api.GetVacationPlansByPeriodHandler(db))
	vacations.Get("/promotions", api.GetPromotionsHandler(db)) //촉진현황 가져오기
}

func registerGroups(apiRouter fiber.Router, db *database.Database) {
	group := apiRouter.Group("groups/:groupID")
	group.Get("/", api.GetGroupHandler(db))
	group.Post("/", api.UpdateGroupHandler(db))
	group.Delete("/", api.DeleteGroupHandler(db))

	members := group.Group("/members")
	members.Get("/", api.GetGroupMembersHandler(db))
	group.Put("/members", api.UpdateGroupMembersHandler(db))

	vacations := group.Group("/vacations")
	vacations.Get("/", api.GetVacationsByPeriodHandler(db))
	vacations.Get("/plans", api.GetVacationPlansByPeriodHandler(db))
}

func registerMembers(apiRouter fiber.Router, db *database.Database) {
	member := apiRouter.Group("/members/:memberID")
	member.Get("/profile", api.GetMemberProfileHandler(db))
	member.Post("/deactivate", api.DeactivateMemberHandler(db))
	member.Delete("/", api.DeleteMemberHandler(db))

	vacations := member.Group("/vacations")
	vacations.Get("/", api.GetVacationsByPeriodHandler(db))
	vacations.Post("/plans", api.CreateVacationPlanHandler(db))
	vacations.Get("/plans", api.GetVacationPlansByPeriodHandler(db))

	notifications := member.Group("/notifications")
	notifications.Get("/", api.GetAllNotificationsHandler(db))
	notifications.Get("/new", api.GetNewNotificationsHandler(db))
	notifications.Post("/:notificationID/approve", api.ApproveNotificationHandler(db))
}

func registerVacations(apiRouter fiber.Router, db *database.Database) {
	vacations := apiRouter.Group("/vacations")

	plan := vacations.Group("/plans/:planId")
	plan.Get("/", api.GetVacationPlanHandler(db))
	plan.Post("/approve", api.ApproveVacationPlanHandler(db))
	plan.Post("/reject", api.RejectVacationPlanHandler(db))
	plan.Post("/cancel-reject", api.CancelRejectVacationPlanHandler(db))
	plan.Post("/", api.UpdateVacationPlanHandler(db))
	plan.Delete("/", api.DeleteVacationPlanHandler(db))

	vacation := vacations.Group("/:vacationID")
	vacation.Get("/", api.GetVacationHandler(db))
	vacation.Post("/", api.UpdateVacationHandler(db))
	vacation.Delete("/", api.DeleteVacationHandler(db))
	vacation.Post("/reject", api.RejectVacationHandler(db))
	vacation.Post("/cancel-reject", api.CancelRejectVacationHandler(db))
}
