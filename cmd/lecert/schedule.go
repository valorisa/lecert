package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/valorisa/lecert/internal/scheduler"
)

func scheduleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule",
		Short: "Manage auto-renewal scheduling",
	}

	cmd.AddCommand(scheduleInstallCmd())
	cmd.AddCommand(scheduleUninstallCmd())
	cmd.AddCommand(scheduleStatusCmd())

	return cmd
}

func scheduleInstallCmd() *cobra.Command {
	var interval string

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install auto-renewal schedule (cron or Task Scheduler)",
		RunE: func(cmd *cobra.Command, args []string) error {
			binary, _ := os.Executable()
			cfg := scheduler.Config{
				Binary:   binary,
				Interval: interval,
			}

			if err := scheduler.Install(cfg); err != nil {
				return fmt.Errorf("install schedule: %w", err)
			}

			fmt.Printf("Auto-renewal scheduled (%s)\n", interval)
			return nil
		},
	}

	cmd.Flags().StringVar(&interval, "interval", "daily", "check interval: hourly, twice-daily, daily")
	return cmd
}

func scheduleUninstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Remove auto-renewal schedule",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := scheduler.Uninstall(); err != nil {
				return fmt.Errorf("uninstall schedule: %w", err)
			}
			fmt.Println("Auto-renewal schedule removed.")
			return nil
		},
	}
}

func scheduleStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check auto-renewal schedule status",
		RunE: func(cmd *cobra.Command, args []string) error {
			status, err := scheduler.Status()
			if err != nil {
				return fmt.Errorf("check status: %w", err)
			}
			fmt.Printf("Auto-renewal: %s\n", status)
			return nil
		},
	}
}
