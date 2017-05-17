package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/labstack/gommon/color"
	"github.com/maps90/go-core/log"
	config "github.com/spf13/viper"
)

type Configuration struct {
	Name, Path string
}

func NewConfiguration(name, path string) *Configuration {
	return &Configuration{
		Name: name,
		Path: path,
	}
}

func (c *Configuration) ReadInConfig() error {
	if c == nil {
		return errors.New("missing configuration.")
	}
	config.SetConfigName(c.Name)
	config.AddConfigPath(c.Path)
	if err := config.ReadInConfig(); err != nil {
		log.New(log.InfoLevelLog, err.Error())
		return err
	}
}

func (c *Configuration) Configure(p string) error {
	if len(p) == 0 {
		c.ReadInConfig()
		return nil
	}
	ext := filepath.Ext(p)
	ext = strings.TrimPrefix(ext, ".")
	config.SetConfigType(ext)

	file, err := os.Open(AbsolutePath(p))
	if err != nil {
		return err
	}
	defer file.Close()
	if err := config.ReadConfig(file); err != nil {
		return fmt.Errorf("%s: %s", color.Red("ERROR"), color.Yellow("config files not found."))
	}

	return nil
}

func AbsolutePath(inPath string) string {
	if strings.HasPrefix(inPath, "$HOME") {
		inPath = userHomeDir() + inPath[5:]
	}

	if strings.HasPrefix(inPath, "$") {
		end := strings.Index(inPath, string(os.PathSeparator))
		inPath = os.Getenv(inPath[1:end]) + inPath[end:]
	}

	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}

	p, err := filepath.Abs(inPath)
	if err == nil {
		return filepath.Clean(p)
	}

	return ""
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
