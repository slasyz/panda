package handle

import (
    "net/http"
)

type ServerFieldsStatic struct {
    Root    string
    Indexes bool
}

func HandleStatic(w http.ResponseWriter, r *http.Request, server *ServerFields) {
    custom := server.Custom.(ServerFieldsStatic)
    _ = custom
    return
}
