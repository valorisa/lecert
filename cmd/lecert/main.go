package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	mode    string
	version = "dev"
)

func main() {
	root := &cobra.Command{
		Use:     "lecert",
		Short:   "Let's Encrypt certificate manager for humans",
		Version: version,
	}

	root.PersistentFlags().StringVar(&mode, "mode", "standard", "interaction mode: novice, standard, expert")

	root.AddCommand(certCmd())
	root.AddCommand(scheduleCmd())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func certCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cert",
		Short: "Certificate operations",
	}

	cmd.AddCommand(obtainCmd())
	cmd.AddCommand(obtainExpertCmd())
	cmd.AddCommand(renewCmd())
	cmd.AddCommand(renewAllCmd())
	cmd.AddCommand(revokeCmd())
	cmd.AddCommand(listCmd())

	return cmd
}
