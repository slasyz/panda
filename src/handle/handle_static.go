package handle

import (
    "github.com/slasyz/panda/src/core"
    "io"
    "mime"
    "net/http"
    "os"
    "path/filepath"
)

// ServerFieldsStatic instances store virtualhost config.
type ServerFieldsStatic struct {
    Root    string
    Indexes bool
}

// LSFilePipeline is used for template rendering (each instance represents file).
type LSFilePipeline struct {
    URL      string
    FileName string
}

// LSPipeline is used for template rendering.
type LSPipeline struct {
    RequestedURL string
    FileList     []LSFilePipeline
    Global       GlobalPipeline
}

// Index file name.
const INDEX_FILE = "index.html"

// NewLSPipeline creates template pipeline containing directory listing.
func NewLSPipeline(dir *os.File, r *http.Request) (result LSPipeline) {
    // Create new instance.
    result = LSPipeline{
        RequestedURL: r.RequestURI,
        Global:       DefaultGlobalPipeline,
    }

    // Create array with file names.
    fileNamesList, err := dir.Readdirnames(-1)
    if err != nil {
        return
    }

    // Create directory listing in terms of pipelines.
    fileList := make([]LSFilePipeline, len(fileNamesList), len(fileNamesList))
    for i, fileName := range fileNamesList {
        fileList[i].FileName = fileName
        fileList[i].URL = filepath.Join(r.RequestURI, fileName)
    }

    result.FileList = fileList
    return
}

// Static virtualhost handler.
func HandleStatic(w http.ResponseWriter, r *http.Request, server *ServerFields) (code int) {
    // Virtualhost config.
    custom := server.Custom.(ServerFieldsStatic)

    // FS path to requested file.
    pathToServe := filepath.Join(custom.Root, r.URL.String())

    core.Debug("return content of file \"%s\"", pathToServe)
    file, err := os.Open(pathToServe)
    defer file.Close()

    // Errors handling.
    if err != nil {
        switch {
        case os.IsNotExist(err):
            code = ReturnError(w, http.StatusNotFound)
        case os.IsPermission(err):
            code = ReturnError(w, http.StatusForbidden)
        default:
            code = ReturnError(w, http.StatusInternalServerError)
        }
        return
    }

    // Getting info about requested file/directory.
    stat, err := file.Stat()
    if err != nil {
        code = ReturnError(w, http.StatusInternalServerError)
        return
    }

    switch mode := stat.Mode(); {
    case mode.IsDir():
        // Trying to find index file.
        indexFile, err := os.Open(filepath.Join(pathToServe, INDEX_FILE))

        if err != nil {
            if os.IsNotExist(err) {
                if custom.Indexes {
                    // Output directory listing.
                    tpl, _ := OpenTemplate(LS_TPL)
                    tpl.Execute(w, NewLSPipeline(file, r))
                    code = http.StatusOK
                } else {
                    // Return 403 error.
                    code = ReturnError(w, http.StatusForbidden)
                }
            } else if os.IsPermission(err) {
                // Not enough rights to read index file.
                code = ReturnError(w, http.StatusForbidden)
            } else {
                // Another error related to index file.
                code = ReturnError(w, http.StatusInternalServerError)
            }
        }

        // Output index file.
        io.Copy(w, indexFile)
        code = http.StatusOK
    case mode.IsRegular():
        // Guessing mime type of requested file.
        mimetype := mime.TypeByExtension(filepath.Ext(pathToServe))
        if mimetype == "" {
            mimetype = "application/octet-stream"
        }
        w.Header().Set("Content-Type", mimetype)

        // Output requested file
        io.Copy(w, file)
        code = http.StatusOK
    }
    return
}
