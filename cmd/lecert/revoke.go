package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/valorisa/lecert/internal/acme"
)

func revokeCmd() *cobra.Command {
	var (
		domain  string
		staging bool
		yes     bool
	)

	cmd := &cobra.Command{
		Use:   "revoke",
		Short: "Revoke a certificate",
		RunE: func(cmd *cobra.Command, args []string) error {
			if domain == "" {
				return fmt.Errorf("--domain is required")
			}

			if !yes {
				fmt.Printf("Revoke certificate for %s? This cannot be undone. [y/N]: ", domain)
				reader := bufio.NewReader(os.Stdin)
				answer, _ := reader.ReadString('\n')
				if strings.TrimSpace(strings.ToLower(answer)) != "y" {
					fmt.Println("Aborted.")
					return nil
				}
			}

			if err := acme.Revoke(domain, staging); err != nil {
				return fmt.Errorf("revoke failed: %w", err)
			}

			fmt.Printf("Certificate for %s revoked successfully.\n", domain)
			return nil
		},
	}

	cmd.Flags().StringVar(&domain, "domain", "", "domain to revoke cert for")
	cmd.Flags().BoolVar(&staging, "staging", false, "use staging environment")
	cmd.Flags().BoolVar(&yes, "yes", false, "skip confirmation prompt")

	return cmd
}
