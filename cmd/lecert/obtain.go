package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/valorisa/lecert/internal/acme"
	"github.com/valorisa/lecert/internal/wizard"
)

func obtainCmd() *cobra.Command {
	var (
		domain      string
		email       string
		challenge   string
		dnsProvider string
		staging     bool
	)

	cmd := &cobra.Command{
		Use:   "obtain",
		Short: "Obtain a new certificate",
		RunE: func(cmd *cobra.Command, args []string) error {
			if mode == "novice" {
				return obtainNovice()
			}

			if domain == "" {
				return fmt.Errorf("--domain is required")
			}
			if email == "" {
				return fmt.Errorf("--email is required")
			}

			req := acme.ObtainRequest{
				Domain:      domain,
				Email:       email,
				Challenge:   challenge,
				DNSProvider: dnsProvider,
				Staging:     staging,
			}

			cert, err := acme.Obtain(req)
			if err != nil {
				return fmt.Errorf("obtain failed: %w", err)
			}

			fmt.Printf("Certificate obtained for %s\n", cert.Domain)
			fmt.Printf("  Expires: %s\n", cert.NotAfter.Format("2006-01-02"))
			fmt.Printf("  Stored:  %s\n", cert.Path)
			return nil
		},
	}

	cmd.Flags().StringVar(&domain, "domain", "", "domain name to obtain cert for")
	cmd.Flags().StringVar(&email, "email", "", "ACME account email")
	cmd.Flags().StringVar(&challenge, "challenge", "http-01", "challenge type: http-01, dns-01")
	cmd.Flags().StringVar(&dnsProvider, "dns-provider", "", "DNS provider for dns-01 challenge (cloudflare, route53, digitalocean)")
	cmd.Flags().BoolVar(&staging, "staging", false, "use Let's Encrypt staging environment")

	return cmd
}

func obtainNovice() error {
	answers, err := wizard.Run()
	if err != nil {
		return fmt.Errorf("wizard: %w", err)
	}

	req := acme.ObtainRequest{
		Domain:    answers.Domain,
		Email:     answers.Email,
		Challenge: answers.Challenge,
		Staging:   false,
	}

	cert, err := acme.Obtain(req)
	if err != nil {
		return fmt.Errorf("obtain failed: %w", err)
	}

	fmt.Printf("\nCertificate obtained for %s\n", cert.Domain)
	fmt.Printf("  Expires: %s\n", cert.NotAfter.Format("2006-01-02"))
	fmt.Printf("  Stored:  %s\n", cert.Path)
	return nil
}
