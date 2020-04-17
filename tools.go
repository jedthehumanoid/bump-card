package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
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

func containsString(list []string, s string) bool {
	for _, val := range list {
		if val == s {
			return true
		}
	}
	return false
}
