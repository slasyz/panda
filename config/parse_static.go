package config

import (
    "errors"
    "github.com/slasyz/panda/handle"
)

// parseServerParameterProxy sets value of "static" server
func parseServerParameterStatic(name, sign, value string, custom *handle.ServerFieldsStatic, currentFileName string) (err error) {
    switch name {
    case "Root":
        custom.Root, err = assignStringValue(name, sign, value)
    case "Indexes":
        custom.Indexes, err = assignBooleanValue(name, sign, value)
    default:
        err = errors.New("unknown field " + name)
    }

    return
}
