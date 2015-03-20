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

func HandleProxy(w http.ResponseWriter, r *http.Request, server *ServerFields) {
    custom := server.Custom.(ServerFieldsProxy)
    _ = custom
    return
}
