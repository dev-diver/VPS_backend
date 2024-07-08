package routes

import (
	"cywell.com/vacation-promotion/app/auth"
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
	registerOrganizes(apiRouter, db)
}

func registerAuth(apiRouter fiber.Router, db *database.Database) {
	auth := apiRouter.Group("/auth")
	auth.Post("/register", api.MakeAdminHandler(db))
	auth.Post("/login", api.LoginHandler(db))
	auth.Post("/logout", api.LogoutHandler())
}

func registerCompanies(apiRouter fiber.Router, db *database.Database) {
	companies := apiRouter.Group("/companies", auth.AuthCheckMiddleware)
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

	organizes := company.Group("/organizes")
	organizes.Get("/", api.GetOrganizesHandler(db))
	organize := organizes.Group("/:organizeID")
	organize.Post("/add", api.AddOrganizeHandler(db)) //name
}

func registerGroups(apiRouter fiber.Router, db *database.Database) {

	groups := apiRouter.Group("/groups", auth.AuthCheckMiddleware)
	group := groups.Group("/:groupID")
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

	members := apiRouter.Group("/members", auth.AuthCheckMiddleware)
	member := members.Group("/:memberID")
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

	vacations := apiRouter.Group("/vacations", auth.AuthCheckMiddleware)
	plans := vacations.Group("/plans")
	plans.Get("/", api.GetVacationPlansByPeriodHandler(db)) //approver, year

	plan := plans.Group("/:planId")
	plan.Get("/", api.GetVacationPlanHandler(db))
	plan.Post("/approve", api.ApproveVacationPlanHandler(db))
	plan.Post("/cancel-approve", api.CancelApproveVacationPlanHandler(db))
	plan.Post("/reject", api.RejectVacationPlanHandler(db))
	plan.Post("/cancel-reject", api.CancelRejectVacationPlanHandler(db))
	plan.Patch("/", api.UpdateVacationPlanHandler(db))
	plan.Delete("/", api.DeleteVacationPlanHandler(db))

	vacation := vacations.Group("/:vacationID")
	vacation.Get("/", api.GetVacationHandler(db))
	vacation.Post("/", api.UpdateVacationHandler(db))
	vacation.Delete("/", api.DeleteVacationHandler(db))
	vacation.Post("/reject", api.RejectVacationHandler(db))
	vacation.Post("/cancel-reject", api.CancelRejectVacationHandler(db))
}

func registerOrganizes(apiRouter fiber.Router, db *database.Database) {

	organizes := apiRouter.Group("/organizes", auth.AuthCheckMiddleware)
	organize := organizes.Group("/:organizeID")
	organize.Put("/", api.UpdateOrganizeHandler(db)) //name 바꾸기
	organize.Delete("/", api.DeleteOrganizeHandler(db))

	members := organize.Group("/members")
	members.Post("/", api.UpdateOrganizeMembersHandler(db)) // [id]
}
