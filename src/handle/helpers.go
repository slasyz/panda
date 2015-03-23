package handle

import (
    //"html/template"
    "net/http"
)

type HTTPErrorPipeline struct {
    Code    int
    Message string
    Global  GlobalPipeline
}

func NewError(code int, message string) HTTPErrorPipeline {
    return HTTPErrorPipeline{
        Code:    code,
        Message: message,
        Global:  DefaultGlobalPipeline,
    }
}

func ReturnError(w http.ResponseWriter, code int) int {
    w.WriteHeader(code)
    tpl, _ := OpenTemplate(ERROR_TPL)
    tpl.Execute(w, NewError(code, http.StatusText(code)))
    return code
}
