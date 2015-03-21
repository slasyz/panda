package handle

import (
    "net/http"
    "time"
)

type ServerFieldsProxy struct {
    URL               string
    Redirect          bool
    Headers           []string
    ConnectTimeout    time.Duration
    ClientMaxBodySize int
}

func HandleProxy(w http.ResponseWriter, r *http.Request, server *ServerFields) {
    custom := server.Custom.(ServerFieldsProxy)
    _ = custom
    return
}
