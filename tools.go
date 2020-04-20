package main

import (
	"io/ioutil"
)

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
