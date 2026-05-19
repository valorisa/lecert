package scheduler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const cronID = "# lecert-auto-renew"

func installCron(cfg Config) error {
	binary, err := resolveAbsPath(cfg.Binary)
	if err != nil {
		return fmt.Errorf("resolve binary path: %w", err)
	}

	schedule := cronSchedule(cfg.Interval)
	entry := fmt.Sprintf("%s %s cert renew-all --quiet %s", schedule, binary, cronID)

	existing, err := currentCrontab()
	if err != nil {
		return fmt.Errorf("read crontab: %w", err)
	}

	if strings.Contains(existing, cronID) {
		return fmt.Errorf("lecert auto-renew already installed (use 'lecert schedule uninstall' first)")
	}

	newCrontab := existing + "\n" + entry + "\n"
	return writeCrontab(newCrontab)
}

func uninstallCron() error {
	existing, err := currentCrontab()
	if err != nil {
		return err
	}

	lines := strings.Split(existing, "\n")
	var filtered []string
	for _, line := range lines {
		if !strings.Contains(line, cronID) {
			filtered = append(filtered, line)
		}
	}

	return writeCrontab(strings.Join(filtered, "\n"))
}

func statusCron() (string, error) {
	existing, err := currentCrontab()
	if err != nil {
		return "unknown", err
	}
	if strings.Contains(existing, cronID) {
		return "installed", nil
	}
	return "not installed", nil
}

func currentCrontab() (string, error) {
	out, err := exec.Command("crontab", "-l").Output()
	if err != nil {
		return "", nil
	}
	return string(out), nil
}

func writeCrontab(content string) error {
	cmd := exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(content)
	return cmd.Run()
}

func cronSchedule(interval string) string {
	switch interval {
	case "hourly":
		return "0 * * * *"
	case "twice-daily":
		return "0 0,12 * * *"
	default:
		return "30 2 * * *"
	}
}

func resolveAbsPath(binary string) (string, error) {
	if filepath.IsAbs(binary) {
		return binary, nil
	}
	path, err := exec.LookPath(binary)
	if err != nil {
		wd, _ := os.Getwd()
		candidate := filepath.Join(wd, binary)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		return binary, nil
	}
	return filepath.Abs(path)
}
