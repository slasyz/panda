package parser

import (
    "errors"
    "github.com/slasyz/panda/src/handle"
)

func parseDefaultParameter(name, sign, value, currentFileName string, dpField *handle.DefaultParameters) (err error) {
    switch name {
    case "AccessLog":
        handle.GlobalParameters.AccessLogger, err = assignLoggerValue(name, sign, value, currentFileName)
    case "ErrorLog":
        handle.GlobalParameters.ErrorLogger, err = assignLoggerValue(name, sign, value, currentFileName)
    case "ReadTimeout":
        handle.GlobalParameters.ReadTimeout, err = assignDurationValue(name, sign, value)
    case "WriteTimeout":
        handle.GlobalParameters.WriteTimeout, err = assignDurationValue(name, sign, value)
    }
    return
}

// parseGlobalParameter sets value of a global parameter
func parseGlobalParameter(name, sign, value, currentFileName string) (err error) {
    switch name {
    case "User":
        handle.GlobalParameters.User, err = assignStringValue(name, sign, value)
    case "AccessLog", "ErrorLog", "ReadTimeout", "WriteTimeout":
        err = parseDefaultParameter(name, sign, value, currentFileName, &handle.GlobalParameters.DefaultParameters)
    case "PathToTPL":
        handle.GlobalParameters.PathToTPL, err = assignPathValue(name, sign, value, currentFileName)
    case "ImportTPLsIntoMemory":
        handle.GlobalParameters.ImportTPLsIntoMemory, err = assignBooleanValue(name, sign, value)
    default:
        err = errors.New("unknown field " + name)
    }

    return
}

// parseServerParameter sets value of specified server's parameter
func parseServerParameter(name, sign, value string, server *handle.ServerFields, currentFileName string) (err error) {
    switch name {
    case "Hostname":
        server.Hostname, err = assignStringValue(name, sign, value)
    case "AccessLog", "ErrorLog", "ReadTimeout", "WriteTimeout":
        err = parseDefaultParameter(name, sign, value, currentFileName, &server.DefaultParameters)
    case "Type":
        server.Type, err = assignStringValue(name, sign, value)

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
