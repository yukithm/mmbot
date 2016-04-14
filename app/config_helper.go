package app

import (
	"log"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/kardianos/osext"
	"github.com/mitchellh/go-homedir"
)

func loadConfig(c *cli.Context) (*Config, error) {
	file := c.String("conf")
	if file == "" {
		file = findConfigFile(c.App.Name)
		if file == "" {
			return DefaultConfig(), nil
		}
	}

	config, err := LoadConfigFile(file)
	if err != nil {
		return nil, err
		log.Fatal(err)
	}

	return config, nil
}

func findConfigFile(appName string) string {
	filename := appName + ".yml"

	// current directory
	if dir, err := os.Getwd(); err == nil {
		if file := findConfigFileInDir(dir, "config", filename); file != "" {
			return file
		}
	}

	// executable directory
	if dir, err := osext.ExecutableFolder(); err == nil {
		if file := findConfigFileInDir(dir, "config", filename); file != "" {
			return file
		}
	}

	// home directory
	// TODO: support XDG_CONFIG_HOME and XDG_CONFIG_DIRS
	if dir, err := homedir.Dir(); err == nil {
		subdir := filepath.Join(".config", appName)
		if file := findConfigFileInDir(dir, subdir, filename); file != "" {
			return file
		}
	}

	return ""
}

func findConfigFileInDir(dir, subdir, filename string) string {
	file := filepath.Join(dir, filename)
	if fileExists(file) {
		return file
	}

	// under sub-directory
	if subdir != "" {
		file = filepath.Join(dir, subdir, filename)
		if fileExists(file) {
			return file
		}
	}

	return ""
}

func fileExists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		return false
	}
	return true
}
