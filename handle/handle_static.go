package handle

import (
    "net/http"
)

type ServerFieldsStatic struct {
    Root    string
    Indexes bool
}

func HandleStatic(request http.Request, server ServerFields) (response http.Response) {
    serverSettings := server.Custom.(ServerFieldsStatic)
    _ = serverSettings
    return
}
