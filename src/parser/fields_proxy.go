package parser

import (
    "errors"
    "github.com/slasyz/panda/src/handle"
    "net/url"
)

// parseServerParameterProxy sets value of "proxy" server instance.
func parseServerParameterProxy(name, sign, value string, custom *handle.ServerFieldsProxy, currentFileName string) (err error) {
    switch name {
    case "URL":
        var rawurl string
        err = assignStringValue(name, sign, value, &rawurl)
        if err != nil {
            return
        }

        url, err := url.Parse(rawurl)
        if err != nil {
            return err
        }

        custom.Scheme = url.Scheme
        custom.Host = url.Host

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
