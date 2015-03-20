package main

import (
    "flag"
    "fmt"
    "github.com/slasyz/panda/src/config"
    "github.com/slasyz/panda/src/core"
    //"github.com/slasyz/panda/src/handle"
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
    errs := config.ParseConfig(*configFile)
    if errs != nil {
        for _, err := range errs {
            core.Log("%s", err)
        }
    }
    if *testConfig {
        core.Log("Config test successful!")
        return
    }

    // TODO
}
