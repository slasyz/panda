package handle

import (
    "errors"
    "github.com/slasyz/panda/src/core"
    "html/template"
    "io/ioutil"
    "os"
    "path/filepath"
    "runtime"
)

type GlobalPipeline struct {
    OS      string
    Arch    string
    Version string
}

const (
    ERROR_TPL = "error.html"
    LS_TPL    = "ls.html"
)

// templates variable contains templates imported to memory.
var templates map[string]*template.Template
var DefaultGlobalPipeline GlobalPipeline

// readTemplateFromFile returns new template.Template instance (from file).
func readTemplateFromFile(tplName string) (tpl *template.Template, err error) {
    tplPath := filepath.Join(GlobalParameters.PathToTPL, tplName)
    tplFile, err := os.Open(tplPath)
    defer tplFile.Close()

    if err != nil {
        return
    }

    stat, _ := tplFile.Stat()
    if stat.IsDir() {
        err = errors.New("file " + tplPath + " should be a template file, not a directory")
        return
    }

    templateBytes, err := ioutil.ReadAll(tplFile)
    tpl, err = template.New(tplName).Parse(string(templateBytes))
    return
}

// Import templates to memory.
func ImportTemplates() (err error) {
    templateFiles := [...]string{ERROR_TPL, LS_TPL}
    templates = make(map[string]*template.Template)

    for _, tplName := range templateFiles {
        core.Debug("import template %s into memory", tplName)
        templates[tplName], err = readTemplateFromFile(tplName)
        if err != nil {
            return
        }
    }

    return
}

// OpenTemplate returns new template.Template instance (from memory or file).
func OpenTemplate(tplName string) (tpl *template.Template, err error) {
    if GlobalParameters.ImportTPLsIntoMemory {
        tpl = templates[tplName]
    } else {
        tpl, err = readTemplateFromFile(tplName)
    }
    return
}

// SetDefaultGlobalPipeline creates GlobalPipeline instance.
func SetDefaultGlobalPipeline() {
    DefaultGlobalPipeline = GlobalPipeline{
        OS:      runtime.GOOS,
        Arch:    runtime.GOARCH,
        Version: core.VERSION,
    }
    return
}
