package config

import (
    "bufio"
    "errors"
    "fmt"
    "github.com/slasyz/panda/core"
    "github.com/slasyz/panda/handle"
    "io"
    "os"
    "path/filepath"
    "regexp"
    "strconv"
)

func parseConfigFile(file io.Reader, currentFileName string, currentServer *int) (errs []configError) {
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
        case virtualHostRegexp.MatchString(currentLine):
            submatches := getRegexpSubmatches(virtualHostRegexp, []string{"ip", "port"}, currentLine)
            ip := submatches[0]
            port, _ := strconv.Atoi(submatches[1])

            server := handle.ServerFields{IP: ip, Port: port}
            handle.GlobalParameters.Servers = append(handle.GlobalParameters.Servers, server)
            *currentServer++
            core.Debug("%s:%d: found new virtualhost with IP %s and port %s", currentFileName, currentLineNumber, ip, port)
        case fieldRegexp.MatchString(currentLine):
            submatches := getRegexpSubmatches(fieldRegexp, []string{"name", "sign", "value"}, currentLine)
            name := submatches[0]
            sign := submatches[1]
            value := submatches[2]

            core.Debug("%s:%d: found field %s with value %s", currentFileName, currentLineNumber, name, value)

            var err error
            if *currentServer == -1 {
                err = parseGlobalParameter(name, sign, value, currentFileName)
            } else {
                err = parseServerParameter(name, sign, value, &handle.GlobalParameters.Servers[*currentServer], currentFileName)
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

                core.Debug("%s:%d: enter the file %s", currentFileName, currentLineNumber, includedFileName)
                errs = append(errs, parseConfigFile(includedFile, includedFileName, currentServer)...)
            } else if command == "include_dir" {
                // TODO
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
    virtualHostRegexp = regexp.MustCompile(`^` +
        `\[VirtualHost\s*\"` +
        `(?P<ip>[^"]+)` +
        `:` +
        `(?P<port>[^":]+)` +
        `\"\]` +
        `$`)
    fieldRegexp = regexp.MustCompile(`^` +
        `(?P<name>[a-zA-Z0-9_]+)` +
        `\s*` +
        `(?P<sign>\+?=)` +
        `\s*` +
        `(?P<value>\"[^"]*\"|\d+(B|KB|MB|GB)?|true|false)` + // string, integer or bool fields
        `$`)
    commandRegexp = regexp.MustCompile(`^` +
        `(?P<command>\w+)` +
        `\(` +
        `(?P<argument>[^)]*)` +
        `\)$`)
    integerRegexp = regexp.MustCompile(`^\d+$`)
    sizeRegexp = regexp.MustCompile(`^(?P<value>\d+)(?P<unit>B|KB|MB|GB)`)

    currentServer := -1
    configErrors := parseConfigFile(file, fileName, &currentServer)
    if len(configErrors) != 0 {
        errs = make([]error, len(configErrors), len(configErrors))

        for i, err := range configErrors {
            errs[i] = errors.New(err.Error())
        }
    }
    return
}
