package parser

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shebang-go/fsmocker/file"
)

var regexPathOpts *regexp.Regexp = regexp.MustCompile(`^(?P<path>.*)\((?P<tags>.*)\)$`)
var regexTag *regexp.Regexp = regexp.MustCompile(`^(?P<tag>err|data|isdir)=(?P<value>.*)`)
var regexFiles *regexp.Regexp = regexp.MustCompile(`^.*\[(?P<files>.*)\]$`)

var regexFilename *regexp.Regexp = regexp.MustCompile(`^([\s\w\.-\:\?_]+)((\(|\[).*)`)

// type Parser struct {
// 	elements []string
// }

func splitString(v string, sep string) []string {
	if strings.HasPrefix(v, string(os.PathSeparator)) {
		return strings.Split(v, sep)[1:]
	}
	retval := strings.Split(v, sep)
	return retval
}

// Parse returns a list of test files
func Parse(v string) []*file.FileInfo {

	retval := make([]*file.FileInfo, 0)
	curPath := "/"
	elements := splitString(v, string(os.PathSeparator))
	for _, v := range elements {

		fname := parseFilename(v)

		tags := parseTags(v)
		newPath := filepath.Join(curPath, fname)
		retval = append(retval, &file.FileInfo{FName: fname, FIsDir: tags.FIsDir, Error: tags.Error, Data: tags.Data, Path: newPath})

		files := parseFiles(v, newPath)
		for _, fi := range files {
			retval = append(retval, fi)
		}
		curPath = newPath
	}
	return retval
}

func parseFilename(v string) string {
	fname := v
	match := regexFilename.FindStringSubmatch(v)
	if match != nil && len(match) >= 0 {
		fname = match[1]
	}
	return strings.TrimSpace(fname)
}

func parseTag(v string) (string, string) {

	match := regexTag.FindStringSubmatch(v)
	if len(match) == 3 {
		if match[1] == "err" || match[1] == "data" || match[1] == "isdir" {
			return match[1], match[2]
		}
	}

	return "", ""
}
func parseTags(v string) file.FileInfo {

	fi := file.FileInfo{FIsDir: true}
	match := regexPathOpts.FindStringSubmatch(v)
	if len(match) == 3 {
		rawTags := splitString(match[2], ",")
		for _, v := range rawTags {
			key, value := parseTag(strings.TrimSpace(v))

			switch key {
			case "err":
				fi.Error = errors.New(value)
			case "data":
				fi.Data = []byte(value)
			case "isdir":
				if value == "false" {
					fi.FIsDir = false
				}
			}
		}
		return fi
	}
	return fi
}

func parseFiles(input string, base string) []*file.FileInfo {

	retval := make([]*file.FileInfo, 0)
	match := regexFiles.FindStringSubmatch(input)
	if len(match) == 2 {
		rawFiles := splitString(match[1], ",")
		for _, v := range rawFiles {
			fname := parseFilename(v)
			tags := parseTags(v)
			retval = append(retval, &file.FileInfo{FName: fname, FIsDir: false, Error: tags.Error, Data: tags.Data, Path: filepath.Join(base, fname)})
		}
	}
	return retval
}
