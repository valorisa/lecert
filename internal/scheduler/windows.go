package scheduler

import (
	"fmt"
	"os/exec"
	"strings"
)

const taskName = "LecertAutoRenew"

func installTaskScheduler(cfg Config) error {
	binary, err := resolveAbsPath(cfg.Binary)
	if err != nil {
		return fmt.Errorf("resolve binary path: %w", err)
	}

	schedule := windowsSchedule(cfg.Interval)

	cmd := exec.Command("schtasks", "/Create",
		"/TN", taskName,
		"/TR", fmt.Sprintf(`"%s" cert renew-all --quiet`, binary),
		"/SC", schedule,
		"/ST", "02:30",
		"/F",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("schtasks create: %w\n%s", err, string(output))
	}
	return nil
}

func uninstallTaskScheduler() error {
	cmd := exec.Command("schtasks", "/Delete", "/TN", taskName, "/F")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("schtasks delete: %w\n%s", err, string(output))
	}
	return nil
}

func statusTaskScheduler() (string, error) {
	cmd := exec.Command("schtasks", "/Query", "/TN", taskName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "does not exist") ||
			strings.Contains(string(output), "n'existe pas") {
			return "not installed", nil
		}
		return "unknown", err
	}
	return "installed", nil
}

func windowsSchedule(interval string) string {
	switch interval {
	case "hourly":
		return "HOURLY"
	case "twice-daily":
		return "DAILY"
	default:
		return "DAILY"
	}
}
