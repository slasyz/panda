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
        fmt.Println("panda version:", core.VERSION)
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

    if handle.GlobalParameters.ImportTPLsIntoMemory {
        err := handle.ImportTemplates()
        if err != nil {
            core.Log("import templates error: %s", err)
            return
        }
    }
    handle.SetDefaultGlobalPipeline()

    // Create http.Server instances
    for i, srv := range handle.GlobalParameters.Servers {
        for _, host := range srv.Hostnames {
            if host == "*" {
                host = ""
            }
            for _, port := range srv.Ports {
                addr := fmt.Sprintf("%s:%d", host, port)
                s := &http.Server{
                    Addr:         addr,
                    Handler:      &handle.GlobalParameters.Servers[i],
                    ReadTimeout:  srv.DefaultParameters.ReadTimeout,
                    WriteTimeout: srv.DefaultParameters.WriteTimeout,
                }
                core.Log("Starting listening to %s", addr)
                go s.ListenAndServe()
            }
        }
    }

    var input string
    fmt.Scanln(&input)
    fmt.Println("done")

    // TODO
}
