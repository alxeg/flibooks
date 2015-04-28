package utils

import (
    "regexp"
    "strings"
    "unicode"
    "unicode/utf8"
)

const (
    RE_NOT_SEPARATORS = `[^[\s\.,:\*\+;\?\\\-—_\(\)\[\]{}<>'"#«»№\/!]+`
)

var (
    re = regexp.MustCompile(RE_NOT_SEPARATORS)
)

func UpperInitial(str string) string {
    if len(str) > 0 {
        process := strings.ToLower(str)
        r, size := utf8.DecodeRuneInString(process)
        return string(unicode.ToUpper(r)) + process[size:]
    }
    return ""
}

func UpperInitialAll(src string) string {
    return re.ReplaceAllStringFunc(src, func(str string) string {
        return UpperInitial(str)
    })
}

func SplitBySeparators(src string) []string {
    return re.FindAllString(src, -1)
}
