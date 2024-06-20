package api

import (
	"cywell.com/vacation-promotion/app/models"
	"cywell.com/vacation-promotion/database"
	"github.com/gofiber/fiber/v2"
)

func GetAllNotificationsHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		memberID := c.Params("memberID")
		var notifications []models.Notification
		if err := db.DB.Where("member_id = ?", memberID).Find(&notifications).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(notifications)
	}
}

func GetNewNotificationsHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		memberID := c.Params("memberID")
		var notifications []models.Notification
		if err := db.DB.Where("member_id = ? AND is_approve = ?", memberID, false).Find(&notifications).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(notifications)
	}
}

func ApproveNotificationHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		notificationID := c.Params("notificationID")
		var notification models.Notification
		if err := db.DB.First(&notification, notificationID).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		//notification.IsApprove = true
		if err := db.DB.Save(&notification).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(notification)
	}
}
