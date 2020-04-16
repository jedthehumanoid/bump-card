package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

// FindFiles is like find
func findFiles(path string) []string {
	files := []string{}
	filepath.Walk(path,
		func(file string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !f.IsDir() {
				files = append(files, file)
			}
			return nil
		})
	return files
}

func readDir(dir string) []string {
	ret := []string{}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		ret = append(ret, f.Name())
	}
	return ret
}

var yamlregexp = regexp.MustCompile("(?ms)^---($.*)^---$")
var tomlregexp = regexp.MustCompile("(?ms)^\\+\\+\\+($.*)^\\+\\+\\+$")

func getFrontmatter(b []byte) (string, []byte, []byte) {
	match := yamlregexp.FindSubmatch(b)
	if len(match) > 1 {
		return "yaml", match[0], match[1]
	}
	match = tomlregexp.FindSubmatch(b)
	if len(match) > 1 {
		return "toml", match[0], match[1]
	}

	return "", []byte{}, []byte{}
}

func containsString(list []string, s string) bool {
	for _, val := range list {
		if val == s {
			return true
		}
	}
	return false
}
