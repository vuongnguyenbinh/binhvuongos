package handler

import "github.com/gofiber/fiber/v2"

// APISuccess returns a standard success JSON response
func APISuccess(data interface{}) fiber.Map {
	return fiber.Map{
		"success": true,
		"data":    data,
		"error":   nil,
	}
}

// APISuccessWithMeta returns success with pagination metadata
func APISuccessWithMeta(data interface{}, page, total int64) fiber.Map {
	return fiber.Map{
		"success": true,
		"data":    data,
		"error":   nil,
		"meta":    fiber.Map{"page": page, "total": total},
	}
}

// APIError returns a standard error JSON response
func APIError(code string, message string) fiber.Map {
	return fiber.Map{
		"success": false,
		"data":    nil,
		"error":   fiber.Map{"code": code, "message": message},
	}
}
