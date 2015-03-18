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

func HandleProxy(request http.Request, server ServerFields) (response http.Response) {
    serverSettings := server.Custom.(ServerFieldsProxy)
    _ = serverSettings
    return
}
