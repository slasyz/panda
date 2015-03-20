package config

import (
    "errors"
    "github.com/slasyz/panda/src/core"
    "log"
    "path/filepath"
    "regexp"
    "strconv"
)

// Line regexps

var lineContentRegexp *regexp.Regexp
var virtualHostRegexp *regexp.Regexp
var fieldRegexp *regexp.Regexp
var commandRegexp *regexp.Regexp

// Value types regexps

var integerRegexp *regexp.Regexp
var sizeRegexp *regexp.Regexp

// Some helper functions

func mustBe(field, typ string) (err error) {
    return errors.New("wrong type of field " + field + " (must be " + typ + ")")
}

func checkAssignSign(name, sign string) (err error) {
    if sign != "=" {
        err = errors.New("you can only assign to this type (must be \"=\" instead of \"+=\")")
    }
    return
}

func checkAppendSign(name, sign string) (err error) {
    if sign != "+=" {
        err = errors.New("you can only append to this type (must be \"+=\" instead of \"=\")")
    }
    return
}

// Integer fields

func assignIntegerValue(name, sign, value string) (result int, err error) {
    err = checkAssignSign(name, sign)

    result, err = parseIntegerValue(name, sign, value)
    return
}

func appendIntegerValue(name, sign, value string, array []int) (result []int, err error) {
    err = checkAppendSign(name, sign)

    resultValue, err := parseIntegerValue(name, sign, value)
    result = append(array, resultValue)
    return
}

func parseIntegerValue(name, sign, value string) (result int, err error) {
    if integerRegexp.MatchString(value) {
        result, _ = strconv.Atoi(value)
    } else {
        err = mustBe(name, "an integer")
    }
    return
}

// Size fields

func assignSizeValue(name, sign, value string) (result int, err error) {
    err = checkAssignSign(name, sign)

    result, err = parseSizeValue(name, sign, value)
    return
}

func parseSizeValue(name, sign, value string) (result int, err error) {
    if sizeRegexp.MatchString(value) {
        submatches := getRegexpSubmatches(sizeRegexp, []string{"value", "unit"}, value)
        result, _ = strconv.Atoi(submatches[0])
        unit := submatches[1]
        factor := 1

        switch unit {
        case "GB":
            factor *= 1024
            fallthrough
        case "MB":
            factor *= 1024
            fallthrough
        case "KB":
            factor *= 1024
        }

        result = factor * result
    } else {
        err = mustBe(name, "data size, e.g. 1000B, 100KB, 10MB, 1GB")
    }
    return
}

// String fields

func assignStringValue(name, sign, value string) (result string, err error) {
    err = checkAssignSign(name, sign)

    result, err = parseStringValue(name, sign, value)
    return
}

func appendStringValue(name, sign, value string, array []string) (result []string, err error) {
    err = checkAppendSign(name, sign)

    resultValue, err := parseStringValue(name, sign, value)
    result = append(array, resultValue)
    return
}

func parseStringValue(name, sign, value string) (result string, err error) {
    if value[0] == '"' && value[len(value)-1] == '"' {
        result = value[1 : len(value)-1] // remove the quotes
    } else {
        err = mustBe(name, "a string")
    }
    return
}

// Path fields

func assignPathValue(name, sign, value, currentFile string) (result string, err error) {
    err = checkAssignSign(name, sign)

    result, err = parsePathValue(name, sign, value, currentFile)
    return
}

func appendPathValue(name, sign, value, currentFile string, array []string) (result []string, err error) {
    err = checkAppendSign(name, sign)

    resultValue, err := parsePathValue(name, sign, value, currentFile)
    result = append(array, resultValue)
    return
}

func parsePathValue(name, sign, value, currentFile string) (result string, err error) {
    if value[0] == '"' && value[len(value)-1] == '"' {
        result := value[1 : len(value)-1] // remove the quotes
        result, _ = filepath.Abs(filepath.Join(filepath.Dir(currentFile), result))
    } else {
        err = mustBe(name, "a string")
    }
    return
}

// Logger fields

func assignLoggerValue(name, sign, value, currentFile string) (result *log.Logger, err error) {
    err = checkAssignSign(name, sign)

    fileName, err := parsePathValue(name, sign, value, currentFile)
    if err != nil {
        return
    }
    result, err = core.OpenLogFile(fileName)

    return
}

// Boolean fields

func assignBooleanValue(name, sign, value string) (result bool, err error) {
    err = checkAssignSign(name, sign)

    if value == "true" {
        result = true
    } else if value == "false" {
        result = false
    } else {
        err = mustBe(name, "true or false")
    }
    return
}
