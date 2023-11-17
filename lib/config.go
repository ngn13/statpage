package lib

import (
	"encoding/json"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type Rule struct {
  // rules for a service 
  Contains      string    `json:"contains"`
  NotContains   string    `json:"not_contains"`
}

type Service struct {
  // results are not set in the config,
  // they are loaded/dumped to the data file
  Results   []Result  `json:"results"`

  // these are also not set in the config
  // just stores the last connection time 
  // and the success status for quick access
  LastTime    int     `json:"last_time"`
  LastSuccess bool    `json:"last_success"`

  Name      string    `json:"name"`
  Address   string    `json:"address"`
  Type      string    `json:"type"`
  Data      string    `json:"data"`
  Link      string    `json:"link"`
  Rule      Rule      `json:"rule"` 
}

type Config struct {
  // all the config options
  Title     string    `json:"title"`
  Address   string    `json:"address"`
  Contact   string    `json:"contact"`
  Interval  string    `json:"interval"`
  Services  []Service `json:"services"`
}

func GetConfig() Config{
  content, err := os.ReadFile(path.Join("config", "config.json"))
  if err != nil {
    log.Fatalf("Error reading config: %s", err)
  }

  var cfg Config
  err = json.Unmarshal(content, &cfg)
  if err != nil {
    log.Fatalf("Error parsing config: %s", err)
  }

  // results are not supposed to be set 
  // from the config, to prevent this
  // we can override the results
  for _, s := range cfg.Services {
    s.Results = []Result{}
  }

  return cfg
}

func GetInterval() time.Duration {
  ival := GetConfig().Interval

  ival_clean := strings.TrimSuffix(ival, string(ival[len(ival)-1])) 
  num, err := strconv.Atoi(ival_clean)
  if err != nil {
    log.Fatalf("Error parsing interval: %s", err)
  }

  if strings.HasSuffix(ival, "m") {
    return time.Duration(num) * (time.Second*60)
  } else if strings.HasSuffix(ival, "s") {
    return time.Duration(num) * time.Second
  } else if strings.HasSuffix(ival, "h") {
    return time.Duration(num) * ((time.Second*60)*60)
  }

  log.Fatal("Bad time format for the interval")
  return -1
}

func GetServices() []Service{
  cfg := GetConfig()
  return cfg.Services
}

