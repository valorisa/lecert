package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/valorisa/lecert/internal/acme"
	"github.com/valorisa/lecert/internal/store"
)

func renewAllCmd() *cobra.Command {
	var (
		quiet   bool
		staging bool
		force   bool
	)

	cmd := &cobra.Command{
		Use:   "renew-all",
		Short: "Renew all certificates due for renewal (expiring within 30 days)",
		RunE: func(cmd *cobra.Command, args []string) error {
			certs, err := store.List()
			if err != nil {
				return fmt.Errorf("list certs: %w", err)
			}

			threshold := 30 * 24 * time.Hour
			var renewed, skipped, failed int

			for _, meta := range certs {
				remaining := time.Until(meta.NotAfter)
				if !force && remaining > threshold {
					skipped++
					continue
				}

				_, err := acme.Renew(meta.Domain, force, staging)
				if err != nil {
					failed++
					if !quiet {
						fmt.Fprintf(cmd.ErrOrStderr(), "FAIL %s: %v\n", meta.Domain, err)
					}
					continue
				}

				renewed++
				if !quiet {
					fmt.Printf("OK   %s renewed\n", meta.Domain)
				}
			}

			if !quiet {
				fmt.Printf("\nSummary: %d renewed, %d skipped, %d failed\n", renewed, skipped, failed)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&quiet, "quiet", false, "suppress output (for scheduled runs)")
	cmd.Flags().BoolVar(&staging, "staging", false, "use staging environment")
	cmd.Flags().BoolVar(&force, "force", false, "force renewal regardless of expiry")

	return cmd
}
