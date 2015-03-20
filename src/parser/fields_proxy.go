package parser

import (
    "errors"
    "github.com/slasyz/panda/src/handle"
)

// parseServerParameterProxy sets value of "proxy" server
func parseServerParameterProxy(name, sign, value string, custom *handle.ServerFieldsProxy, currentFileName string) (err error) {
    switch name {
    case "URL":
        custom.URL, err = assignStringValue(name, sign, value)
    case "Redirect":
        custom.Redirect, err = assignBooleanValue(name, sign, value)
    case "Headers":
        custom.Headers, err = appendStringValue(name, sign, value, custom.Headers)
    case "ConnectTimeout":
        custom.ConnectTimeout, err = assignIntegerValue(name, sign, value)
    case "ClientMaxBodySize":
        custom.ClientMaxBodySize, err = assignSizeValue(name, sign, value)
    default:
        err = errors.New("unknown field " + name)
    }

    return
}
