package handle

import (
    //"html/template"
    "net/http"
)

// HTTPErrorPipeline contains information about error.
type HTTPErrorPipeline struct {
    Code    int
    Message string
    Global  GlobalPipeline
}

// Returns template pipeline containing error info.
func NewError(code int, message string) HTTPErrorPipeline {
    return HTTPErrorPipeline{
        Code:    code,
        Message: message,
        Global:  DefaultGlobalPipeline,
    }
}

// ReturnError writes error page into http.ResponseWriter.
func ReturnError(w http.ResponseWriter, code int) int {
    w.WriteHeader(code)
    tpl, _ := OpenTemplate(ERROR_TPL)
    tpl.Execute(w, NewError(code, http.StatusText(code)))
    return code
}
