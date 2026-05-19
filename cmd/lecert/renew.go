package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/valorisa/lecert/internal/acme"
)

func renewCmd() *cobra.Command {
	var (
		domain  string
		force   bool
		staging bool
	)

	cmd := &cobra.Command{
		Use:   "renew",
		Short: "Renew an existing certificate",
		RunE: func(cmd *cobra.Command, args []string) error {
			if domain == "" {
				return fmt.Errorf("--domain is required")
			}

			cert, err := acme.Renew(domain, force, staging)
			if err != nil {
				return fmt.Errorf("renew failed: %w", err)
			}

			fmt.Printf("Certificate renewed for %s\n", cert.Domain)
			fmt.Printf("  New expiry: %s\n", cert.NotAfter.Format("2006-01-02"))
			return nil
		},
	}

	cmd.Flags().StringVar(&domain, "domain", "", "domain to renew")
	cmd.Flags().BoolVar(&force, "force", false, "force renewal even if not near expiry")
	cmd.Flags().BoolVar(&staging, "staging", false, "use staging environment")

	return cmd
}
