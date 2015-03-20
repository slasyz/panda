package handle

import (
    //"github.com/slasyz/panda/core"
    "log"
    "net/http"
)

// HandleFunc is handler function for HTTP request type
type HandleFunc func(w http.ResponseWriter, r *http.Request)

// GlobalParameters is variable containing all panda parameters
var GlobalParameters struct {
    Listen               []int
    User                 string
    AccessLogger         *log.Logger
    ErrorLogger          *log.Logger
    PathToTPL            string
    ImportTPLsIntoMemory bool
    Servers              []ServerFields
}

// Server is type containing each server parameters
type ServerFields struct {
    IP       string
    Port     int
    Hostname string
    Type     string
    Custom   interface{}
    Handler  HandleFunc
}
