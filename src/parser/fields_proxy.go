package parser

import (
    "errors"
    "github.com/slasyz/panda/src/handle"
)

// parseServerParameterProxy sets value of "proxy" server
func parseServerParameterProxy(name, sign, value string, custom *handle.ServerFieldsProxy, currentFileName string) (err error) {
    switch name {
    case "URL":
        err = assignStringValue(name, sign, value, &custom.URL)
    case "Redirect":
        err = assignBooleanValue(name, sign, value, &custom.Redirect)
    case "Headers":
        err = appendStringValue(name, sign, value, &custom.Headers)
    case "ConnectTimeout":
        err = assignDurationValue(name, sign, value, &custom.ConnectTimeout)
    case "ClientMaxBodySize":
        err = assignSizeValue(name, sign, value, &custom.ClientMaxBodySize)
    default:
        err = errors.New("unknown field " + name)
    }

    return
}
