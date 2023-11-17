package lib

import (
	"time"

	"github.com/gofiber/fiber/v2/log"
)

var check_ticker = time.NewTicker(GetInterval())
var result_ticker = time.NewTicker(time.Duration(24)*time.Hour)
var Checkchan = make(chan struct{})
var Checktime time.Time 

func checkServices() {
  log.Info("Checking service status")
  services := GetConfig().Services
  for _, s := range services {
    ok, t := CheckService(s)
    log.Infof("Service: %s | Type: %s | Status: %t | Time: %dms", s.Name, s.Type, ok, t)
    SaveResult(s, t, ok)
  }
  Checktime = time.Now()
}

func addResult() {
  log.Info("Resetting results")
  services := Services
  
  for i := range services {
    Services[i].Results = append(Services[i].Results, Result{
      SuccessRate: 100,
      TimeAverage: 0,
      Success: []bool{},
      Time: []int{},
    })
    
    if len(Services[i].Results) >= 20 {
      copy(Services[i].Results[0:], Services[i].Results[1:]) 
      Services[i].Results[len(Services[i].Results)-1] = Result{}     
      Services[i].Results = Services[i].Results[:len(Services[i].Results)-1]
    }
  }

  SaveResults()
}

func Loop() {
  checkServices()

  for {
    select {
      case <- result_ticker.C:
        addResult()
        checkServices()

      case <- check_ticker.C:
        checkServices()

      case <- Checkchan:
        result_ticker.Stop()
        check_ticker.Stop()
        return
    }
  }
} 

