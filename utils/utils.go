package utils

import (
    "regexp"
    "strings"
    "unicode"
    "unicode/utf8"
)

const (
    RE_NOT_SEPARATORS    = `[^[\s\.,:\*\+;\?\\\-—_\(\)\[\]{}<>'"#«»№\/!]+`
    RE_UNSUPPORTED_CHARS = `[\\\/\*\+\?]`
)

var (
    re_split = regexp.MustCompile(RE_NOT_SEPARATORS)
    re_unsup = regexp.MustCompile(RE_UNSUPPORTED_CHARS)
)

func ReplaceUnsupported(str string) string {
    return re_unsup.ReplaceAllString(str, "_")
}

func UpperInitial(str string) string {
    if len(str) > 0 {
        process := strings.ToLower(str)
        r, size := utf8.DecodeRuneInString(process)
        return string(unicode.ToUpper(r)) + process[size:]
    }
    return ""
}

func UpperInitialAll(src string) string {
    return re_split.ReplaceAllStringFunc(src, func(str string) string {
        return UpperInitial(str)
    })
}

func SplitBySeparators(src string) []string {
    return re_split.FindAllString(src, -1)
}
