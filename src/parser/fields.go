package parser

import (
    "errors"
    "github.com/slasyz/panda/src/handle"
)

// parseGlobalParameter sets value of a global parameter
func parseGlobalParameter(name, sign, value, currentFile string) (err error) {
    switch name {
    case "User":
        handle.GlobalParameters.User, err = assignStringValue(name, sign, value)
    case "Listen":
        handle.GlobalParameters.Listen, err = appendIntegerValue(name, sign, value, handle.GlobalParameters.Listen)
    case "AccessLog":
        handle.GlobalParameters.AccessLogger, err = assignLoggerValue(name, sign, value, currentFile)
    case "ErrorLog":
        handle.GlobalParameters.ErrorLogger, err = assignLoggerValue(name, sign, value, currentFile)
    case "PathToTPL":
        handle.GlobalParameters.PathToTPL, err = assignPathValue(name, sign, value, currentFile)
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
    case "Type":
        server.Type, err = assignStringValue(name, sign, value)

        switch server.Type {
        case "static":
            server.Custom = handle.ServerFieldsStatic{}
            server.Handler = handle.HandleStatic
        case "proxy":
            server.Custom = handle.ServerFieldsProxy{}
            server.Handler = handle.HandleProxy
        }
    default:
        switch server.Type {
        case "static":
            srvStatic := server.Custom.(handle.ServerFieldsStatic)
            err = parseServerParameterStatic(name, sign, value, &srvStatic, currentFileName)
            server.Custom = srvStatic
        case "proxy":
            srvProxy := server.Custom.(handle.ServerFieldsProxy)
            err = parseServerParameterProxy(name, sign, value, &srvProxy, currentFileName)
            server.Custom = srvProxy
        }
    }

    return
}
