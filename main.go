package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/ngn13/statpage/lib"
)

func main(){
  lib.LoadResults()
  go lib.Loop()

  engine := html.New("./views", ".html") 
  app := fiber.New(fiber.Config{
    Views: engine,
  })

  app.Use(logger.New())
  app.Static("/", "./public")

  app.Get("/", func(c *fiber.Ctx) error {
    t := time.Since(lib.Checktime)
    chtime := fmt.Sprintf(
      "%ds ago", 
      int(t.Seconds()),
    )

    if t.Minutes() > 1 {
      chtime = fmt.Sprintf(
        "%dm and %ds ago", 
        int(t.Minutes()), int(t.Seconds())-(int(t.Minutes())*60),
      )
    }

    if t.Hours() > 1 {
      chtime = fmt.Sprintf("%dh and %dm ago", 
        int(t.Hours()),
        int(t.Minutes())-(int(t.Hours())*60), 
      )
    }

    return c.Render("index", fiber.Map{
      "title": lib.GetConfig().Title,
      "contact": lib.GetConfig().Contact,
      "services": lib.Services,
      "checktime": chtime,
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

  addr := lib.GetConfig().Address
  log.Infof("Starting the application on %s", addr)
  log.Fatal(app.Listen(addr))
  log.Info("Stopped the application")
  close(lib.Checkchan)
}
