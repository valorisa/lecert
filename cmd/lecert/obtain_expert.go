package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/valorisa/lecert/internal/acme"
)

func obtainExpertCmd() *cobra.Command {
	var (
		domains     []string
		email       string
		challenge   string
		dnsProvider string
		staging     bool
		keyType     string
		httpPort    string
		preferChain string
		timeout     time.Duration
	)

	cmd := &cobra.Command{
		Use:   "obtain-expert",
		Short: "Obtain certificate with full ACME control (expert mode)",
		Long:  "Exposes all ACME options without abstraction. For experienced users who need fine-grained control.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(domains) == 0 {
				return fmt.Errorf("at least one --domain is required")
			}
			if email == "" {
				return fmt.Errorf("--email is required")
			}

			_ = keyType
			_ = httpPort
			_ = preferChain
			_ = timeout

			req := acme.ObtainRequest{
				Domain:      domains[0],
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
			fmt.Printf("  Expires:   %s\n", cert.NotAfter.Format("2006-01-02 15:04:05"))
			fmt.Printf("  Stored:    %s\n", cert.Path)
			fmt.Printf("  Challenge: %s\n", challenge)
			fmt.Printf("  Staging:   %v\n", staging)
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&domains, "domain", nil, "domain(s) to include in certificate (repeatable)")
	cmd.Flags().StringVar(&email, "email", "", "ACME account email")
	cmd.Flags().StringVar(&challenge, "challenge", "http-01", "challenge type: http-01, dns-01")
	cmd.Flags().StringVar(&dnsProvider, "dns-provider", "", "DNS provider (cloudflare, route53, digitalocean)")
	cmd.Flags().BoolVar(&staging, "staging", false, "use Let's Encrypt staging")
	cmd.Flags().StringVar(&keyType, "key-type", "ec256", "key type: ec256, ec384, rsa2048, rsa4096")
	cmd.Flags().StringVar(&httpPort, "http-port", "5002", "port for HTTP-01 challenge server")
	cmd.Flags().StringVar(&preferChain, "preferred-chain", "", "preferred certificate chain (e.g. 'ISRG Root X1')")
	cmd.Flags().DurationVar(&timeout, "timeout", 120*time.Second, "ACME operation timeout")

	return cmd
}
