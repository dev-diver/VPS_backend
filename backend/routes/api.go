package routes

import (
	"cywell.com/vacation-promotion/app/controllers/api"
	"cywell.com/vacation-promotion/database"
	"github.com/gofiber/fiber/v2"
)

func RegisterAPI(apiRouter fiber.Router, db *database.Database) {
	registerAuth(apiRouter, db)
	registerCompanies(apiRouter, db)
	registerMembers(apiRouter, db)
	registerGroups(apiRouter, db)
}

func registerAuth(apiRouter fiber.Router, db *database.Database) {
	auth := apiRouter.Group("/auth")
	auth.Post("/login", api.LoginHandler(db))
	auth.Post("/logout", api.LogoutHandler(db))
}

func registerCompanies(apiRouter fiber.Router, db *database.Database) {
	companies := apiRouter.Group("/companies")
	companies.Post("/", api.CreateCompanyHandler(db))
	companies.Delete("/", api.DeleteCompaniesHandler(db))

	company := companies.Group("/:companyID")
	company.Get("/", api.GetCompanyHandler(db))
	company.Post("/", api.UpdateCompanyHandler(db))

	vacations := company.Group("/vacations")
	vacations.Get("/period/:year/:month?", api.GetVacationsByYearMonthHandler(db))
	vacations.Post("/:vacationID", api.CreateVacationHandler(db))
	vacations.Put("/:vacationID", api.UpdateVacationHandler(db))
	vacations.Delete("/:vacationID", api.DeleteVacationHandler(db))
	vacations.Post("/:vacationID/promotion", api.PromoteVacationHandler(db))
	vacations.Get("/promotions", api.GetPromotionsHandler(db))

	members := company.Group("/members")
	members.Get("/", api.GetMembersHandler(db))
	members.Get("/search", api.SearchMembersHandler(db)) // keyword
	members.Post("/", api.CreateMembersHandler(db))      // []
	member := members.Group("/:memberID")
	member.Get("/profile", api.GetMemberProfileHandler(db))
	member.Post("/deactivate", api.DeactivateMemberHandler(db))
	member.Delete("/", api.DeleteMemberHandler(db))

	groups := company.Group("/groups")
	groups.Get("/", api.GetGroupsHandler(db))
	groups.Post("/", api.CreateGroupHandler(db))
	group := groups.Group("/:groupID")
	group.Post("/", api.UpdateGroupHandler(db))
	group.Delete("/", api.DeleteGroupHandler(db))
	group.Put("/members", api.UpdateGroupMembersHandler(db)) // member key
}

func registerMembers(apiRouter fiber.Router, db *database.Database) {
	member := apiRouter.Group("/members/:memberID")
	vacations := member.Group("/vacations")
	vacations.Post("/", api.ApplyVacationHandler(db))
	vacation := vacations.Group("/:vacationID")
	vacation.Post("/", api.UpdateVacationHandler(db))
	vacation.Delete("/", api.CancelVacationHandler(db))
	vacation.Post("/approve", api.ApproveVacationHandler(db))
	vacation.Post("/reject", api.RejectVacationHandler(db))
	vacations.Get("/period/:year/:month?", api.GetMemberVacationsHandler(db))

	notifications := member.Group("/notifications")
	notifications.Get("/", api.GetAllNotificationsHandler(db))
	notifications.Get("/new", api.GetNewNotificationsHandler(db))
	notifications.Post("/:notificationID/approve", api.ApproveNotificationHandler(db))
}

func registerGroups(apiRouter fiber.Router, db *database.Database) {
	groups := apiRouter.Group("/groups")
	groups.Post("/", api.CreateGroupHandler(db))
	group := groups.Group("/:groupID")
	group.Get("/", api.GetGroupHandler(db))
	group.Put("/", api.UpdateGroupHandler(db))
	members := group.Group("/members")
	members.Post("/", api.AddGroupMembersHandler(db)) // memberId[]
	member := members.Group("/:memberID")
	member.Delete("/", api.DeleteGroupMemberHandler(db))
	group.Delete("/", api.DeleteGroupHandler(db))
}
