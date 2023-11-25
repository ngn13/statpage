package lib

import (
	"time"

	"github.com/gofiber/fiber/v2/log"
)

var check_ticker = time.NewTicker(GetInterval())
var Checkchan = make(chan struct{})
var LastChecked time.Time 

func checkServices() {
  log.Info("Checking service status")
  services := GetConfig().Services
  for _, s := range services {
    ok, t := CheckService(s)
    log.Infof("Service: %s | Type: %s | Status: %t | Time: %dms", s.Name, s.Type, ok, t)
    SaveResult(s, t, ok)
  }
  LastChecked = time.Now()
}

func Loop() {
  checkServices()

  for {
    select {
      case <- check_ticker.C:
        checkServices()

      case <- Checkchan:
        check_ticker.Stop()
        return
    }
  }
} 

