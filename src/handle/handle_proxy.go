package handle

import (
    "github.com/slasyz/panda/src/core"
    "net/http"
    "net/http/httputil"
    "strings"
    "time"
)

// ServerFieldsProxy instances store virtualhost config.
type ServerFieldsProxy struct {
    Scheme string
    Host   string

    Redirect          bool
    Headers           []string
    ConnectTimeout    time.Duration
    ClientMaxBodySize int
}

// Splits "Host: $host"-like string into two parts and replaces all variables
func parseHeaderFromRaw(header string) (key, value string) {
    arr := strings.Split(header, ":")
    key = strings.TrimSpace(arr[0])
    value = strings.TrimSpace(arr[1])

    core.Debug("Added header %s: %s", key, value)

    return
}

// Proxy virtualhost handler.
func HandleProxy(w http.ResponseWriter, r *http.Request, server *ServerFields) (code int) {
    custom := server.Custom.(ServerFieldsProxy)
    _ = custom

    director := func(req *http.Request) {
        req = r
        req.URL.Scheme = custom.Scheme
        req.URL.Host = custom.Host

        for _, header := range custom.Headers {
            key, value := parseHeaderFromRaw(header)
            req.Header.Add(key, value)
        }
    }

    //core.Debug("REQ URL: %s", r.URL)

    proxy := &httputil.ReverseProxy{Director: director}
    proxy.ServeHTTP(w, r)

    return
}
