package main

import (
    "flag"
    "fmt"
    "github.com/slasyz/panda/src/core"
    "github.com/slasyz/panda/src/handle"
    "github.com/slasyz/panda/src/parser"
    "log"
    "net/http"
)

func main() {
    configFile := flag.String("config", "conf/main.conf", "path to config file")
    testConfig := flag.Bool("test-config", false, "test config and exit")
    core.DebugFlag = flag.Bool("d", false, "very detailed output level")
    version := flag.Bool("v", false, "show version and exit")
    flag.Parse()

    // show version and exit
    if *version {
        fmt.Println("panda version: 0.0.1")
        return
    }

    // config parsing
    core.OpenedLoggers = make(map[string]*log.Logger)
    defer core.CloseLogFiles()
    errs := parser.ParseConfig(*configFile)
    if errs != nil {
        for _, err := range errs {
            core.Log("%s", err)
        }
    }
    if *testConfig {
        core.Log("Config test successful!")
        return
    }

    // Begin to listen
    for i, srv := range handle.GlobalParameters.Servers {
        // create new http.Server
        s := &http.Server{
            Addr:         srv.GetAddr(),
            Handler:      srv,
            ReadTimeout:  srv.DefaultParameters.ReadTimeout,
            WriteTimeout: srv.DefaultParameters.WriteTimeout,
        }
        core.Log("Starting listening to %s (server #%d)", srv.GetAddr(), i+1)
        go s.ListenAndServe()
    }

    var input string
    fmt.Scanln(&input)
    fmt.Println("done")

    // TODO
}
