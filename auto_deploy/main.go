package main

import (
	"os/exec"

	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	app.Get("/invest_indicator", func(c *fiber.Ctx) error {
		cmd := exec.Command("./deploy_invest_indicator.sh")
		if err := cmd.Run(); err != nil {
			return c.SendString(err.Error())
		}
		return nil
	})
	app.Get("/lolche_bot", func(c *fiber.Ctx) error {
		cmd := exec.Command("./deploy_lolche_bot.sh")
		if err := cmd.Run(); err != nil {
			return c.SendString(err.Error())
		}
		return nil
	})
	app.Listen(":10000")
}