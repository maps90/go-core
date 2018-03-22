package core

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/labstack/gommon/color"
	config "github.com/spf13/viper"

	_ "github.com/spf13/viper/remote" // required for remote access
)

type Configuration struct {
	Name, Path, Remote, URL string
}

func NewConfiguration(name, path, remote, url string) *Configuration {
	return &Configuration{
		Name:   name,
		Remote: remote,
		URL:    url,
		Path:   path,
	}
}

func (c *Configuration) ReadInConfig() error {
	// find consul environment variables
	if c.URL == "" {
		color.Println(color.Green(fmt.Sprintf("⇨ using local config: '%v.yaml'", c.Name)))
		return setLocalConfig(c.Name)
	}

	color.Print(color.Green(fmt.Sprintf("⇨ connecting to consul(%v) ... ", url)))
	err := setRemoteConfig(c.URL, c.Remote)
	if err != nil {
		color.Println(color.Red("failed"))
	} else {
		color.Println(color.Green("success!"))
	}

	return nil
}

func setLocalConfig(conf string) (err error) {
	// set file type
	config.SetConfigType("yaml")
	config.SetConfigName(conf)
	config.AddConfigPath(".")

	err = config.ReadInConfig()
	return
}

func setRemoteConfig(url string, conf string) (err error) {
	err = config.AddRemoteProvider("consul", url, conf)
	if err != nil {
		return err
	}

	config.SetConfigType("yaml")
	// read from remote config.
	if err := config.ReadRemoteConfig(); err != nil {
		return err
	}

	return nil
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
