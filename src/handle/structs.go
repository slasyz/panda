package handle

import (
    //"github.com/slasyz/panda/core"
    "fmt"
    "log"
    "net"
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
    IP         string
    Port       int
    Hostname   string
    Type       string
    Custom     interface{}
    HandleFunc func(http.ResponseWriter, *http.Request, *ServerFields)
}

func (server ServerFields) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    server.HandleFunc(w, r, &server)
}

func (server ServerFields) GetAddr() string {
    portString := fmt.Sprintf("%d", server.Port)
    fmt.Println(">>>>>> ", net.JoinHostPort(server.IP, portString))
    return net.JoinHostPort(server.IP, portString)
}

func NewServer(ip string, port int) (server *ServerFields) {
    server = &ServerFields{IP: ip,
        Port:              port,
        DefaultParameters: GlobalParameters.DefaultParameters,
    }
    return
}
