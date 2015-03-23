package handle

import (
    "github.com/slasyz/panda/src/core"
    "io"
    "net/http"
    "os"
    "path/filepath"
)

type ServerFieldsStatic struct {
    Root    string
    Indexes bool
}

type LSFilePipeline struct {
    URL      string
    FileName string
}

type LSPipeline struct {
    RequestedURL string
    FileList     []LSFilePipeline
    Global       GlobalPipeline
}

func NewLSPipeline(dir *os.File, r *http.Request) (result LSPipeline) {
    result = LSPipeline{
        RequestedURL: r.RequestURI,
        Global:       DefaultGlobalPipeline,
    }

    fileNamesList, err := dir.Readdirnames(-1)
    if err != nil {
        return
    }

    fileList := make([]LSFilePipeline, len(fileNamesList), len(fileNamesList))
    for i, fileName := range fileNamesList {
        fileList[i].FileName = fileName
        fileList[i].URL = filepath.Join(r.RequestURI, fileName)
    }

    result.FileList = fileList
    return
}

const INDEX_FILE = "index.html"

func HandleStatic(w http.ResponseWriter, r *http.Request, server *ServerFields) {
    custom := server.Custom.(ServerFieldsStatic)
    pathToServe := filepath.Join(custom.Root, r.URL.String())

    core.Debug("return content of file \"%s\"", pathToServe)
    file, err := os.Open(pathToServe)
    defer file.Close()

    if err != nil {
        switch {
        case os.IsNotExist(err):
            ReturnError(w, http.StatusNotFound)
        case os.IsPermission(err):
            ReturnError(w, http.StatusForbidden)
        default:
            ReturnError(w, http.StatusInternalServerError)
        }
        return
    }

    stat, err := file.Stat()
    if err != nil {
        ReturnError(w, http.StatusInternalServerError)
        return
    }

    switch mode := stat.Mode(); {
    case mode.IsDir():
        indexFile, err := os.Open(filepath.Join(pathToServe, INDEX_FILE))
        if err != nil {
            if os.IsNotExist(err) {
                if custom.Indexes {
                    tpl, _ := OpenTemplate(LS_TPL)
                    tpl.Execute(w, NewLSPipeline(file, r))
                } else {
                    ReturnError(w, http.StatusForbidden)
                }
            } else {
                ReturnError(w, http.StatusInternalServerError)
            }
        }
        io.Copy(w, indexFile)
    case mode.IsRegular():
        io.Copy(w, file)
    }
    return
}
