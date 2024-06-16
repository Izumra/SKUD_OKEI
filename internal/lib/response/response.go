package response

import "github.com/gofiber/fiber/v3"

type Body struct {
	Data  any   `json:"data,omitempty"`
	Error error `json:"error,omitempty"`
}

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
