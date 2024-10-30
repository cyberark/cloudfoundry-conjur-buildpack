package main

import (
	"io/ioutil"
	"os"
	"strings"
)

const devShmPath = "/dev/shm"

type tempFactory struct {
	path  string
	files []string
}

// Create a new temporary file factory.
// defer cleanup() if you want the files removed.
func newTempFactory() tempFactory {
	return tempFactory{path: defaultTempPath()}
}

// Default temporary file path
// Returns /dev/shm if it is a directory, otherwise home dir of current user
// Else returns the system default
func defaultTempPath() string {
	fi, err := os.Stat(devShmPath)
	if err == nil && fi.Mode().IsDir() {
		return devShmPath
	}
	home, err := os.UserHomeDir()
	if err == nil {
		dir, _ := ioutil.TempDir(home, ".tmp")
		return dir
	}
	return os.TempDir()
}

// Create a temp file with given value. Returns the path.
func (tf *tempFactory) push(bytes []byte) string {
	f, _ := os.CreateTemp(tf.path, ".conjur-env")
	defer f.Close()

	f.Write(bytes)
	name := f.Name()
	tf.files = append(tf.files, name)
	return name
}

// Remove the temporary files created with this factory.
func (tf *tempFactory) cleanup() {
	for _, file := range tf.files {
		os.Remove(file)
	}
	// Also remove the tempdir if it's not /dev/shm
	if !strings.Contains(tf.path, devShmPath) {
		os.Remove(tf.path)
	}
	tf = nil
}
