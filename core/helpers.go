package core

import (
    "errors"
    "fmt"
    "log"
    "os"
    "time"
)

var DebugFlag *bool
var OpenedFiles map[string]*log.Logger

func Debug(str string, a ...interface{}) {
    if *DebugFlag {
        Log("DEBUG: "+str, a...)
    }
}

func Log(str string, a ...interface{}) {
    str = fmt.Sprintf("[%s] %s\n", time.Now().Format("2006/02/01 15:04:05"), str)
    fmt.Printf(str, a...)
}

func OpenLogFile(fileName string) (logger *log.Logger, err error) {
    logger, ok := OpenedFiles[fileName]
    if ok {
        return
    } else {
        file, err := os.Open(fileName)
        if err != nil {
            err = errors.New("cannot open file " + fileName)
        }
        logger = log.New(file, "", log.LstdFlags)
    }
    return
}
