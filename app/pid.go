package app

import (
	"io/ioutil"
	"os"
	"strconv"
)

type PIDFile string

func NewPIDFile(pidfile string) (PIDFile, error) {
	if pidfile == "" {
		return "", nil
	}

	pidfile, err := absPath(pidfile)
	if err != nil {
		return "", err
	}

	return PIDFile(pidfile), nil
}

func (p PIDFile) Exists() bool {
	if p == "" {
		return false
	}
	return fileExists(string(p))
}

func (p PIDFile) Create() error {
	if p == "" {
		return nil
	}

	pid := strconv.Itoa(os.Getpid())
	return ioutil.WriteFile(string(p), []byte(pid), 0644)
}

func (p PIDFile) Remove() error {
	err := os.Remove(string(p))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
