package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type Result struct {
  SuccessRate int     `json:"success_rate"` 
  TimeAverage int     `json:"time_average"`
  Success     []bool  `json:"success"`
  Time        []int   `json:"time"`
  Date        string  `json:"date"`
}

// data is just a json list containing all the services 
// as well as check results
type Data struct {
  Services []Service `json:"services"`
}

var Services []Service

func CheckDataDir() {
  st, err := os.Stat("data")

  // if we got an error and its NOT a NotExist error
  // this means we cannot access the directory,
  // this may be caused by permissions
  // in this case we just return an error
  if err != nil && !os.IsNotExist(err) {
    log.Fatalf("Error checking data directory: %s", err) 
  }

  // if we got an error and its a NotExist error
  // this means the directory does not exists
  // in this case we can just create it
  if err != nil && os.IsNotExist(err) {
    err = os.Mkdir("data", os.ModePerm)
    if err != nil {
      log.Fatalf("Error creating data directory: %s", err)
    }
    return
  }

  // if there is no errors but the stat tells us that 
  // data is a file, then we return an error
  if err == nil && !st.IsDir() {
    log.Fatal("Cannot access data directory because it's a file")
  }
}

func LoadResults() {
  CheckDataDir()

  content, err := os.ReadFile(path.Join("data", "data.json"))
  if err != nil && !os.IsNotExist(err) {
    log.Fatalf("Error reading data file: %s", err)
  }

  if err != nil && os.IsNotExist(err) {
    return
  }

  // creating a temp data object
  var data Data
  // we can then use this temp object to 
  // load all the services 

  err = json.Unmarshal(content, &data)
  if err != nil {
    log.Fatalf("Error parsing data file: %s", err)
  }

  Services = data.Services
}

func SaveResults() {
  CheckDataDir()
  
  // create a temp data object
  var data Data
  // this time we will use this to 
  // dump all services
  data.Services = Services

  content, err := json.Marshal(&data)
  if err != nil {
    log.Fatalf("Error dumping data file: %s", err)
  }

  err = os.WriteFile(path.Join("data", "data.json"), content, os.ModePerm)
  if err != nil {
    log.Fatalf("Error writing to data file: %s", err)
  }
}

func CalcRate(result Result) Result {
  // really ugly rate calculations
  result.SuccessRate = 100
  persuccess := 100/len(result.Success)
  for _, s := range result.Success {
    if !s {
      result.SuccessRate = result.SuccessRate - persuccess
    }
  }

  result.TimeAverage = 0
  for _, t := range result.Time {
    result.TimeAverage += t
  }
  result.TimeAverage = result.TimeAverage / len(result.Time)

  return result
}

func CheckRemoved() {
  cfg_services := GetConfig().Services
  for i, s := range Services {
    found := false
    for _, c := range cfg_services {
      if s.Address == c.Address && s.Type == c.Type {
        found = true
        break
      }
    }

    if !found {
      copy(Services[i:], Services[i+1:]) 
      Services[len(Services)-1] = Service{}  
      Services = Services[:len(Services)-1]
    }
  }

  SaveResults()
}

func SaveResult(service Service, t int, ok bool) {
  CheckRemoved()

  found := false
  now := time.Now()
  date := fmt.Sprintf(
    "%d/%d/%d", 
    now.Day(), now.Month(), now.Year(),
  )

  for i := range Services {
    if Services[i].Address == service.Address && Services[i].Type == service.Type {
      indx := len(Services[i].Results)-1
      Services[i].Results[indx].Time = append(Services[i].Results[indx].Time, t)
      Services[i].Results[indx].Success = append(Services[i].Results[indx].Success, ok)
      Services[i].Results[indx].Date = date

      Services[i].LastTime = t 
      Services[i].LastSuccess = ok 
      found = true
    }
  }

  if !found {
    service.Results = append(service.Results, Result{
      SuccessRate: 100,
      TimeAverage: 0,
      Success: []bool{ok},
      Time: []int{t},
      Date: date,
    })

    service.LastTime = t
    service.LastSuccess = ok 
    Services = append(Services, service)
  }

  for _, s := range Services {
    if s.Address == service.Address && s.Type == service.Type { 
      last := len(s.Results)-1
      s.Results[last] = CalcRate(s.Results[last])
    } 
  }

  SaveResults()
}
