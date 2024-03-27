package response

import "github.com/gofiber/fiber/v3"

func BadRes(err error) fiber.Map {
	return fiber.Map{
		"data":  nil,
		"error": err.Error(),
	}
}

func SuccessRes(data any) fiber.Map {
	return fiber.Map{
		"data":  data,
		"error": nil,
	}
}
