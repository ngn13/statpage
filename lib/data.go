package lib

import (
	"encoding/json"
	"os"
	"path"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type Result struct {
  Time        int             `json:"time"`
  Success     bool            `json:"success"`
  Stamp       time.Time       `json:"stamp"`
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

func LoadData() {
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

func SaveData() {
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

  SaveData()
}

func SaveResult(service Service, t int, ok bool) {
  CheckRemoved()

  for i := range Services {
    if Services[i].Address != service.Address || Services[i].Type != service.Type {
      continue
    }
    
    Services[i].Results = append(Services[i].Results, Result{
      Time: t,
      Success: ok,
      Stamp: time.Now(),
    })

    Services[i].LastTime = t 
    Services[i].LastSuccess = ok  

    if(len(Services[i].Results) > GetConfig().Reset) {
      copy(Services[i].Results[0:], Services[i].Results[1:])
      Services[i].Results[len(Services[i].Results)-1] = Result{}
      Services[i].Results = Services[i].Results[:len(Services[i].Results)-1]
    }

    SaveData()
    return
  }

  service.Results = append(service.Results, Result{
    Success: ok,
    Time: t,
    Stamp: time.Now(),
  })

  service.LastTime = t
  service.LastSuccess = ok 
  Services = append(Services, service)
  
  SaveData()
}
