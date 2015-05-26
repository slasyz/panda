package handle

import (
    "net/http"
    "time"
)

// ServerFieldsProxy instances store virtualhost config.
type ServerFieldsProxy struct {
    URL               string
    Redirect          bool
    Headers           []string
    ConnectTimeout    time.Duration
    ClientMaxBodySize int
}

// Proxy virtualhost handler.
func HandleProxy(w http.ResponseWriter, r *http.Request, server *ServerFields) (code int) {
    custom := server.Custom.(ServerFieldsProxy)
    _ = custom

    return
}
