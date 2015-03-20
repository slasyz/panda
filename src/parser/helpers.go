package parser

import (
    "fmt"
    "regexp"
)

// configError is type containing config error (Captain Obvious at your
// service)
type configError struct {
    fileName   string
    lineNumber int
    message    string
}

func (err *configError) Error() string {
    return fmt.Sprintf("%s:%d: %s", err.fileName, err.lineNumber, err.message)
}

func getRegexpSubmatches(reg *regexp.Regexp, names []string, s string) (submatches []string) {
    index := reg.FindStringSubmatchIndex(s)
    submatches = make([]string, len(names), len(names))
    for i, name := range names {
        submatches[i] = string(reg.ExpandString([]byte{}, "$"+name, s, index))
    }
    return
}
