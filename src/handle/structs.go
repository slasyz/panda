package handle

import (
    //"github.com/slasyz/panda/core"
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
    HandleFunc func(http.ResponseWriter, *http.Request, *ServerFields)
}

func (server *ServerFields) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    server.HandleFunc(w, r, server)
    return
}
