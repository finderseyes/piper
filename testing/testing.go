// Package testing ...
package testing

import (
	"os"
	"path"
	"runtime"
)

// nolint: dogsled,gochecknoinits // disable since it's too simple.
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
