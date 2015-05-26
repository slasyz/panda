package core

import (
    "errors"
    "fmt"
    "log"
    "os"
    "time"
)

var (
    DebugFlag     *bool
    OpenedLoggers map[string]*log.Logger
    OpenedFiles   []*os.File
)

const VERSION = "0.0.1"

// Output debug information.
func Debug(str string, a ...interface{}) {
    if *DebugFlag {
        Log("DEBUG: "+str, a...)
    }
}

// Output general information.
func Log(str string, a ...interface{}) {
    str = fmt.Sprintf("[%s] %s\n", time.Now().Format("2006/02/01 15:04:05"), str)
    fmt.Printf(str, a...)
}

// OpenLogFile opens log file and creates log.Logger instance.
func OpenLogFile(fileName string) (logger *log.Logger, err error) {
    Debug("opening log file %s", fileName)
    logger, ok := OpenedLoggers[fileName]
    if ok {
        return
    } else {
        file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
        if err != nil {
            err = errors.New("cannot open file " + fileName)
        }
        logger = log.New(file, "", log.LstdFlags)

        OpenedLoggers[fileName] = logger
        OpenedFiles = append(OpenedFiles, file)
    }
    return
}

// CloseLogFiles closes all open log files.
func CloseLogFiles() {
    for _, file := range OpenedFiles {
        file.Close()
    }
}
