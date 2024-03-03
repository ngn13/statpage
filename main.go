package main

import (
	"fmt"
	"os"
	"path"
	"time"

  "github.com/gofiber/fiber/v2"
  "github.com/gofiber/fiber/v2/log"
  "github.com/gofiber/fiber/v2/middleware/logger"
  "github.com/gofiber/template/django/v3"
  "github.com/ngn13/statpage/lib"
)

func CheckTimePassed(t time.Time) string {
  diff := time.Since(t)
  res := fmt.Sprintf(
    "%ds ago", 
    int(diff.Seconds()),
  )

  if diff.Minutes() > 1 { 
    res = fmt.Sprintf(
      "%dm %ds ago", 
      int(diff.Minutes()), int(diff.Seconds())-(int(diff.Minutes())*60),
    )
  }

  if diff.Hours() > 1 {
    res = fmt.Sprintf("%dh %dm ago", 
      int(diff.Hours()),
      int(diff.Minutes())-(int(diff.Hours())*60), 
    )
  }

  return res
}

func main(){
  lib.LoadData()
  go lib.Loop()

  engine := django.New("./views", ".html")
  app := fiber.New(fiber.Config{
    DisableStartupMessage: true,
    Views: engine,
  })

  app.Use(logger.New())
  app.Static("/", "./public")

  app.Get("/", func(c *fiber.Ctx) error {
    return c.Render("index", fiber.Map{
      "cfg": lib.GetConfig(),
      "services": lib.Services,
      "lastcheck": lib.LastChecked,
      "checktime": CheckTimePassed,
    })
  })

  app.Get("/theme.css", func(c *fiber.Ctx) error {
    theme, err := os.ReadFile(path.Join("config", "theme.css"))
    if err != nil{
      if !os.IsNotExist(err) {
        log.Errorf("Error reading theme: %s", err)
      }

      return c.Status(404).Send([]byte(""))
    }

    c.Set("Content-Type", "text/css; charset=utf-8")
    return c.Send(theme)
  })

  app.Get("/favicon.ico", func(c *fiber.Ctx) error {
    return c.Status(404).SendString("I don't have an icon :/")
  })

  app.Get("*", func(c *fiber.Ctx) error {
    return c.Redirect("/")
  })

  addr := lib.GetConfig().Address
  log.Infof("Starting the application on %s", addr)
  log.Fatal(app.Listen(addr))
  log.Info("Stopped the application")
  close(lib.Checkchan)
}
