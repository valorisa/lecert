package scheduler

import (
	"fmt"
	"runtime"
)

type Config struct {
	Binary   string
	Interval string
}

func Install(cfg Config) error {
	if cfg.Binary == "" {
		cfg.Binary = "lecert"
	}
	if cfg.Interval == "" {
		cfg.Interval = "daily"
	}

	switch runtime.GOOS {
	case "linux", "darwin":
		return installCron(cfg)
	case "windows":
		return installTaskScheduler(cfg)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func Uninstall() error {
	switch runtime.GOOS {
	case "linux", "darwin":
		return uninstallCron()
	case "windows":
		return uninstallTaskScheduler()
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func Status() (string, error) {
	switch runtime.GOOS {
	case "linux", "darwin":
		return statusCron()
	case "windows":
		return statusTaskScheduler()
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}
