package parser

import (
    "errors"
    "github.com/slasyz/panda/src/handle"
)

// parseServerParameterProxy sets value of "static" server
func parseServerParameterStatic(name, sign, value string, custom *handle.ServerFieldsStatic, currentFileName string) (err error) {
    switch name {
    case "Root":
        err = assignStringValue(name, sign, value, &custom.Root)
    case "Indexes":
        err = assignBooleanValue(name, sign, value, &custom.Indexes)
    default:
        err = errors.New("unknown field " + name)
    }

    return
}
