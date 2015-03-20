package handle

import (
    "net/http"
)

type ServerFieldsStatic struct {
    Root    string
    Indexes bool
}

func (server *ServerFields) HandleStatic(w http.ResponseWriter, r *http.Request) {
    custom := server.Custom.(ServerFieldsStatic)
    _ = custom
    return
}
