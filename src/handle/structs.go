package handle

import (
    "fmt"
    "github.com/slasyz/panda/src/core"
    "log"
    "net/http"
    "time"
)

type DefaultParameters struct {
    AccessLogger *log.Logger
    ErrorLogger  *log.Logger
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
}

// GlobalParameters is variable containing all panda parameters
var GlobalParameters struct {
    DefaultParameters
    User                 string
    PathToTPL            string
    ImportTPLsIntoMemory bool
    Servers              []ServerFields
}

// Server is type containing each server parameters
type ServerFields struct {
    DefaultParameters
    Hostnames  []string
    Ports      []int
    Type       string
    Custom     interface{}
    HandleFunc func(http.ResponseWriter, *http.Request, *ServerFields) int
}

func (server *ServerFields) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    code := server.HandleFunc(w, r, server)

    timeStr := time.Now().Format("02/Jan/2006 15:04:05 -0700")
    logStr := fmt.Sprintf("%s \"%s %s\" %d \"%s\"", r.RemoteAddr, r.Method, r.RequestURI, code, r.UserAgent())

    core.Log(logStr)
    server.AccessLogger.Printf("[%s] %s\n", timeStr, logStr)
    return
}
