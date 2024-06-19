package routes

import (
	"cywell.com/vacation-promotion/database"
	"github.com/gofiber/fiber/v2"
)

func RegisterAPI(api fiber.Router, db *database.Database) {
	registerAuth(api, db)
	registerCompanies(api, db)
	registerMembers(api, db)
	registerGroups(api, db)
}

func registerAuth(api fiber.Router, db *database.Database) {
	auth := api.Group("/auth")
	auth.Post("/login", loginHandler)
	auth.Post("/logout", logoutHandler)
}

func registerCompanies(api fiber.Router, db *database.Database) {
	companies := api.Group("/companies")
	companies.Post("/", createCompanyHandler)
	companies.Delete("/", deleteCompaniesHandler)

	company := companies.Group("/:companyID")
	company.Get("/", getCompanyHandler)
	company.Post("/", updateCompanyHandler)

	vacations := company.Group("/vacations")
	vacations.Get("/period/:year/:month?", getVacationsByYearMonthHandler)
	vacations.Post("/:vacationID", createVacationHandler)
	vacations.Put("/:vacationID", updateVacationHandler)
	vacations.Delete("/:vacationID", deleteVacationHandler)
	vacations.Post("/:vacationID/promotion", promoteVacationHandler)
	vacations.Get("/promotions", getPromotionsHandler)

	members := company.Group("/members")
	members.Get("/", getMembersHandler)
	members.Get("/search", searchMembersHandler) //keyword
	members.Post("/", createMembersHandler)      //[]
	member := members.Group("/:memberID")
	member.Get("/profile", getMemberProfileHandler)
	member.Post("/deactivate", deactivateMemberHandler)
	member.Delete("/", deleteMemberHandler)

	groups := company.Group("/groups")
	groups.Get("/", getGroupsHandler)
	groups.Post("/", createGroupHandler)
	group := groups.Group("/:groupID")
	group.Post("/", updateGroupHandler)
	group.Delete("/", deleteGroupHandler)
	group.Put("/members", updateGroupMembersHandler) // member key
}

func registerMembers(api fiber.Router, db *database.Database) {
	member := api.Group("/members/:memberID")
	vacations := member.Group("/vacations")
	vacations.Post("/", applyVacationHandler)
	vacation := vacations.Group("/:vacationID")
	vacation.Post("/", updateVacationHandler)
	vacation.Delete("/", cancelVacationHandler)
	vacation.Post("/approve", approveVacationHandler)
	vacation.Post("/reject", rejectVacationHandler)
	vacations.Get("/period/:year/:month?", getMemberVacationsHandler)

	notifications := member.Group("/notifications")
	notifications.Get("/", getAllNotificationsHandler)
	notifications.Get("/new", getNewNotificationsHandler)
	notifications.Post("/:notificationID/approve", approveNotificationHandler)
}

func registerGroups(api fiber.Router, db *database.Database) {
	groups := api.Group("/groups")
	groups.Post("/", createGroupHandler)
	group := groups.Group("/:groupID")
	group.Get("/", getGroupHandler)
	group.Put("/", updateGroupHandler)
	members := group.Group("/members")
	members.Post("/", addGroupMembersHandler) // memberId[]
	member := members.Group("/:memberID")
	member.Delete("/", deleteGroupMemberHandler)
	group.Delete("/", deleteGroupHandler)
}

// Handler functions (placeholders)
func loginHandler(c *fiber.Ctx) error                   { return c.SendStatus(fiber.StatusOK) }
func logoutHandler(c *fiber.Ctx) error                  { return c.SendStatus(fiber.StatusOK) }
func createCompanyHandler(c *fiber.Ctx) error           { return c.SendStatus(fiber.StatusOK) }
func deleteCompaniesHandler(c *fiber.Ctx) error         { return c.SendStatus(fiber.StatusOK) }
func getCompanyHandler(c *fiber.Ctx) error              { return c.SendStatus(fiber.StatusOK) }
func updateCompanyHandler(c *fiber.Ctx) error           { return c.SendStatus(fiber.StatusOK) }
func getVacationsByYearMonthHandler(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) }
func createVacationHandler(c *fiber.Ctx) error          { return c.SendStatus(fiber.StatusOK) }
func updateVacationHandler(c *fiber.Ctx) error          { return c.SendStatus(fiber.StatusOK) }
func deleteVacationHandler(c *fiber.Ctx) error          { return c.SendStatus(fiber.StatusOK) }
func promoteVacationHandler(c *fiber.Ctx) error         { return c.SendStatus(fiber.StatusOK) }
func getPromotionsHandler(c *fiber.Ctx) error           { return c.SendStatus(fiber.StatusOK) }
func getMembersHandler(c *fiber.Ctx) error              { return c.SendStatus(fiber.StatusOK) }
func searchMembersHandler(c *fiber.Ctx) error           { return c.SendStatus(fiber.StatusOK) }
func createMembersHandler(c *fiber.Ctx) error           { return c.SendStatus(fiber.StatusOK) }
func getMemberProfileHandler(c *fiber.Ctx) error        { return c.SendStatus(fiber.StatusOK) }
func deactivateMemberHandler(c *fiber.Ctx) error        { return c.SendStatus(fiber.StatusOK) }
func deleteMemberHandler(c *fiber.Ctx) error            { return c.SendStatus(fiber.StatusOK) }
func getGroupsHandler(c *fiber.Ctx) error               { return c.SendStatus(fiber.StatusOK) }
func getGroupHandler(c *fiber.Ctx) error                { return c.SendStatus(fiber.StatusOK) }
func createGroupHandler(c *fiber.Ctx) error             { return c.SendStatus(fiber.StatusOK) }
func updateGroupHandler(c *fiber.Ctx) error             { return c.SendStatus(fiber.StatusOK) }
func deleteGroupHandler(c *fiber.Ctx) error             { return c.SendStatus(fiber.StatusOK) }
func updateGroupMembersHandler(c *fiber.Ctx) error      { return c.SendStatus(fiber.StatusOK) }
func applyVacationHandler(c *fiber.Ctx) error           { return c.SendStatus(fiber.StatusOK) }
func cancelVacationHandler(c *fiber.Ctx) error          { return c.SendStatus(fiber.StatusOK) }
func approveVacationHandler(c *fiber.Ctx) error         { return c.SendStatus(fiber.StatusOK) }
func rejectVacationHandler(c *fiber.Ctx) error          { return c.SendStatus(fiber.StatusOK) }
func getMemberVacationsHandler(c *fiber.Ctx) error      { return c.SendStatus(fiber.StatusOK) }
func getAllNotificationsHandler(c *fiber.Ctx) error     { return c.SendStatus(fiber.StatusOK) }
func getNewNotificationsHandler(c *fiber.Ctx) error     { return c.SendStatus(fiber.StatusOK) }
func approveNotificationHandler(c *fiber.Ctx) error     { return c.SendStatus(fiber.StatusOK) }
func addGroupMembersHandler(c *fiber.Ctx) error         { return c.SendStatus(fiber.StatusOK) }
func deleteGroupMemberHandler(c *fiber.Ctx) error       { return c.SendStatus(fiber.StatusOK) }
