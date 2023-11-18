package lib

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

func CheckGET(s Service, url string,) (bool, int) {
  req, err := http.NewRequest("GET", url, nil)
  var start time.Time
  var elapsed int64

  trace := &httptrace.ClientTrace{
    ConnectStart: func(_, _ string){ start = time.Now() },
    GotFirstResponseByte: func(){ elapsed = time.Since(start).Milliseconds() },
  } 

  req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
  res, err := http.DefaultClient.Do(req)

  if err != nil {
    return false, 0
  }

  defer res.Body.Close()
  raw, err := io.ReadAll(res.Body)
  if err != nil {
    log.Errorf("Error reading body from %s: %s", url, err)
    return false, 0
  }

  if s.Rule.NotContains != "" {
    if strings.Contains(string(raw), s.Rule.NotContains) {
      return false, int(elapsed)
    }
  }

  if s.Rule.Contains != "" {
    if !strings.Contains(string(raw), s.Rule.Contains) {
      return false, int(elapsed)
    }
  }

  return true, int(elapsed)

}

func CheckTCP(s Service, useTls bool) (bool, int) {
  var conn net.Conn 
  var err error

  if useTls {
    var tlscfg *tls.Config 
    conn, err = tls.Dial("tcp", s.Address, tlscfg)
  }else {
    conn, err = net.Dial("tcp", s.Address)
  }

  if err != nil {
    return false, 0
  }
  defer conn.Close()

  start := time.Now()
  if s.Data != "" {
    _, err = conn.Write([]byte(s.Data))
    if err != nil {
      log.Errorf("Error sending data to %s: %s", s.Address, err)
      return false, 0
    }
  }

  buff := make([]byte, 1024)
  _, err = conn.Read(buff)
  if err != nil {
    log.Errorf("Error receiving data from %s: %s", s.Address, err)
    return false, 0
  }
  elapsed := time.Since(start).Milliseconds() 

  if s.Rule.NotContains != "" {
    if strings.Contains(string(buff), s.Rule.NotContains) {
      return false, int(elapsed)
    }
  }

  if s.Rule.Contains != "" {
    if !strings.Contains(string(buff), s.Rule.Contains) {
      return false, int(elapsed)
    }
  }

  return true, int(elapsed)
}

func CheckHTTPS(s Service) (bool, int) {
  url := "https://"+s.Address+s.Data
  return CheckGET(s, url)
}

func CheckHTTP(s Service) (bool, int) {
  url := "http://"+s.Address+s.Data
  return CheckGET(s, url)
}

func CheckService(s Service) (bool, int) {
  switch strings.ToLower(s.Type) {
  case "http":
    return CheckHTTP(s)
  case "https":
    return CheckHTTPS(s)
  case "tcp":
    return CheckTCP(s, false)
  case "tls":
    return CheckTCP(s, true)
  }

  log.Fatalf("Unknown service type for %s: %s", s.Address, s.Type)
  return false, 0
} 
