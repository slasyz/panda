package parser

import (
    "errors"
    "github.com/slasyz/panda/src/handle"
)

// parseDefaultParameter check general parameters for any server type.
func parseDefaultParameter(name, sign, value, currentFileName string, dpField *handle.DefaultParameters) (err error) {
    switch name {
    case "AccessLog":
        err = assignLoggerValue(name, sign, value, currentFileName, &handle.GlobalParameters.AccessLogger)
    case "ErrorLog":
        err = assignLoggerValue(name, sign, value, currentFileName, &handle.GlobalParameters.ErrorLogger)
    case "ReadTimeout":
        err = assignDurationValue(name, sign, value, &handle.GlobalParameters.ReadTimeout)
    case "WriteTimeout":
        err = assignDurationValue(name, sign, value, &handle.GlobalParameters.WriteTimeout)
    }
    return
}

// parseGlobalParameter sets value of a global parameter.
func parseGlobalParameter(name, sign, value, currentFileName string) (err error) {
    switch name {
    case "User":
        err = assignStringValue(name, sign, value, &handle.GlobalParameters.User)
    case "AccessLog", "ErrorLog", "ReadTimeout", "WriteTimeout":
        err = parseDefaultParameter(name, sign, value, currentFileName, &handle.GlobalParameters.DefaultParameters)
    case "PathToTPL":
        err = assignPathValue(name, sign, value, currentFileName, &handle.GlobalParameters.PathToTPL)
    case "ImportTPLsIntoMemory":
        err = assignBooleanValue(name, sign, value, &handle.GlobalParameters.ImportTPLsIntoMemory)
    default:
        err = errors.New("unknown field " + name)
    }

    return
}

// parseServerParameter sets value of specified server type parameter.
func parseServerParameter(name, sign, value string, server *handle.ServerFields, currentFileName string) (err error) {
    switch name {
    case "Hostnames":
        err = appendStringValue(name, sign, value, &server.Hostnames)
    case "Ports":
        err = appendIntegerValue(name, sign, value, &server.Ports)
    case "AccessLog", "ErrorLog", "ReadTimeout", "WriteTimeout":
        err = parseDefaultParameter(name, sign, value, currentFileName, &server.DefaultParameters)
    case "Type":
        err = assignStringValue(name, sign, value, &server.Type)

        switch server.Type {
        case "static":
            server.Custom = handle.ServerFieldsStatic{}
            server.HandleFunc = handle.HandleStatic
        case "proxy":
            server.Custom = handle.ServerFieldsProxy{}
            server.HandleFunc = handle.HandleProxy
        }
    default:
        switch server.Type {
        case "static":
            srvCustom := server.Custom.(handle.ServerFieldsStatic)
            err = parseServerParameterStatic(name, sign, value, &srvCustom, currentFileName)
            server.Custom = srvCustom
        case "proxy":
            srvCustom := server.Custom.(handle.ServerFieldsProxy)
            err = parseServerParameterProxy(name, sign, value, &srvCustom, currentFileName)
            server.Custom = srvCustom
        }
    }

    return
}
