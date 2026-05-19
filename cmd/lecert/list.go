package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/valorisa/lecert/internal/store"
)

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List managed certificates",
		RunE: func(cmd *cobra.Command, args []string) error {
			certs, err := store.List()
			if err != nil {
				return fmt.Errorf("list failed: %w", err)
			}

			if len(certs) == 0 {
				fmt.Println("No certificates managed yet.")
				return nil
			}

			fmt.Printf("%-30s %-12s %-10s\n", "DOMAIN", "EXPIRES", "STATUS")
			for _, c := range certs {
				status := "valid"
				if time.Now().After(c.NotAfter) {
					status = "expired"
				} else if time.Until(c.NotAfter) < 30*24*time.Hour {
					status = "renew-soon"
				}
				fmt.Printf("%-30s %-12s %-10s\n", c.Domain, c.NotAfter.Format("2006-01-02"), status)
			}
			return nil
		},
	}
}
