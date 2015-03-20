package handle

import (
    "net/http"
)

type ServerFieldsProxy struct {
    URL               string
    Redirect          bool
    Headers           []string
    ConnectTimeout    int
    ClientMaxBodySize int
}

func (server *ServerFields) HandleProxy(w http.ResponseWriter, r *http.Request) {
    custom := server.Custom.(ServerFieldsProxy)
    _ = custom
    return
}
