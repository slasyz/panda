package parser

import (
    "bufio"
    "errors"
    "fmt"
    "github.com/slasyz/panda/src/core"
    "github.com/slasyz/panda/src/handle"
    "io"
    "os"
    "path/filepath"
    "regexp"
    "sort"
)

// Starting to parse config file.
func parseConfigFile(file io.Reader, currentFileName string) (errs []configError) {
    scanner := bufio.NewScanner(file)
    currentLineNumber := 0

    for scanner.Scan() {
        currentLineNumber++
        currentLine := scanner.Text()
        if !lineContentRegexp.MatchString(currentLine) {
            errs = append(errs, configError{currentFileName, currentLineNumber, "syntax error (unpaired quotation mark?)"})
        }
        currentLine = getRegexpSubmatches(lineContentRegexp, []string{"content"}, currentLine)[0]

        switch {
        case currentLine == "":
            continue
        case currentLine == "[VirtualHost]":
            server := handle.ServerFields{DefaultParameters: handle.GlobalParameters.DefaultParameters}
            handle.GlobalParameters.Servers = append(handle.GlobalParameters.Servers, server)
            core.Debug("%s:%d: found new virtualhost", currentFileName, currentLineNumber)
        case fieldRegexp.MatchString(currentLine):
            submatches := getRegexpSubmatches(fieldRegexp, []string{"name", "sign", "value"}, currentLine)
            name := submatches[0]
            sign := submatches[1]
            value := submatches[2]

            core.Debug("%s:%d: found field %s with value %s", currentFileName, currentLineNumber, name, value)

            var err error
            if len(handle.GlobalParameters.Servers) == 0 {
                err = parseGlobalParameter(name, sign, value, currentFileName)
            } else {
                currentServer := len(handle.GlobalParameters.Servers) - 1
                err = parseServerParameter(name, sign, value, &handle.GlobalParameters.Servers[currentServer], currentFileName)
            }

            if err != nil {
                errs = append(errs, configError{currentFileName, currentLineNumber, err.Error()})
            }
        case commandRegexp.MatchString(currentLine):
            submatches := getRegexpSubmatches(commandRegexp, []string{"command", "argument"}, currentLine)
            command := submatches[0]
            argument := submatches[1]

            if command == "include" {
                includedFileName := filepath.Join(filepath.Dir(currentFileName), argument)
                includedFile, err := os.Open(includedFileName)
                if err != nil {
                    errs = append(errs, configError{currentFileName, currentLineNumber, "file " + includedFileName + " can not be open"})
                }

                // Getting file/directory info.
                stat, err := includedFile.Stat()
                if err != nil {
                    errs = append(errs, configError{currentFileName, currentLineNumber, "file " + includedFileName + " can not be open"})
                }

                switch mode := stat.Mode(); {
                case mode.IsRegular():
                    core.Debug("%s:%d: enter the file %s", currentFileName, currentLineNumber, includedFileName)
                    errs = append(errs, parseConfigFile(includedFile, includedFileName)...)
                case mode.IsDir():
                    errs = append(errs, configError{currentFileName, currentLineNumber, includedFileName + " is not a file"})
                }
            } else if command == "include_dir" {
                includedDirectoryName := filepath.Join(filepath.Dir(currentFileName), argument)
                includedDirectory, err := os.Open(includedDirectoryName)

                // Error while opening file/directory.
                if err != nil {
                    switch {
                    case os.IsNotExist(err):
                        errs = append(errs, configError{currentFileName, currentLineNumber, "directory " + includedDirectoryName + " is not exists"})
                    case os.IsPermission(err):
                        errs = append(errs, configError{currentFileName, currentLineNumber, "directory " + includedDirectoryName + " can not be open (permission error)"})
                    default:
                        errs = append(errs, configError{currentFileName, currentLineNumber, "directory " + includedDirectoryName + " can not be open"})
                    }
                }

                // Getting file/directory info.
                stat, err := includedDirectory.Stat()
                if err != nil {
                    errs = append(errs, configError{currentFileName, currentLineNumber, "directory " + includedDirectoryName + " can not be open"})
                }

                switch mode := stat.Mode(); {
                case mode.IsDir():
                    // Getting directory content.
                    fileNamesList, err := includedDirectory.Readdirnames(-1)
                    if err != nil {
                        errs = append(errs, configError{currentFileName, currentLineNumber, "error getting directory " + includedDirectoryName + " content"})
                    }
                    sort.Strings(fileNamesList)

                    for i, includedFileName := range fileNamesList {
                        fileNamesList[i] = filepath.Join(includedDirectoryName, includedFileName)
                    }

                    // Include all files from the directory.
                    core.Debug("%s:%d: including all files from directory %s", currentFileName, currentLineNumber, includedDirectoryName)
                    for _, includedFileName := range fileNamesList {

                        includedFile, err := os.Open(includedFileName)
                        if err != nil {
                            errs = append(errs, configError{currentFileName, currentLineNumber, "file " + includedFileName + " can not be open"})
                        }

                        core.Debug("%s:%d: enter the file %s", currentFileName, currentLineNumber, includedFileName)
                        errs = append(errs, parseConfigFile(includedFile, includedFileName)...)
                    }

                case mode.IsRegular():
                    errs = append(errs, configError{currentFileName, currentLineNumber, includedDirectoryName + " is not a directory"})
                }
            }
        default:
            errs = append(errs, configError{currentFileName, currentLineNumber, "wrong line format"})
        }
    }
    return
}

func ParseConfig(fileName string) (errs []error) {
    file, err := os.Open(fileName)
    if err != nil {
        return []error{errors.New(fmt.Sprintf("Cannot read config file %s (%s)", fileName, err))}
    }
    defer file.Close()

    // Compile the regexps once
    lineContentRegexp = regexp.MustCompile(`^\s*` +
        `(?P<content>(("[^"]*")|([^"]))*?)` +
        `\s*(\/\/.*)?\s*$`)
    fieldRegexp = regexp.MustCompile(`^` +
        `(?P<name>[a-zA-Z0-9_]+)` +
        `\s*` +
        `(?P<sign>\+?=)` +
        `\s*` +
        `(?P<value>\"[^"]*\"|\d+(B|KB|MB|GB|ns|us|µs|ms|s|m|h)?|true|false)` + // string, integer or bool fields
        `$`)
    commandRegexp = regexp.MustCompile(`^` +
        `(?P<command>\w+)` +
        `\(` +
        `(?P<argument>[^)]*)` +
        `\)$`)
    integerRegexp = regexp.MustCompile(`^\d+$`)
    sizeRegexp = regexp.MustCompile(`^(?P<value>\d+)(?P<unit>B|KB|MB|GB)`)

    configErrors := parseConfigFile(file, fileName)
    if len(configErrors) != 0 {
        errs = make([]error, len(configErrors), len(configErrors))

        for i, err := range configErrors {
            errs[i] = errors.New(err.Error())
        }
    }
    return
}
