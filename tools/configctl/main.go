package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "configctl",
		Short: "Configuration management tool for EduGo API",
		Long: `configctl is a CLI tool to manage configuration variables in the EduGo API Mobile project.

It helps you:
  - Add new configuration variables (public or secret)
  - Validate configuration files
  - Generate documentation

Examples:
  configctl add database.redis.host --type string --default localhost --desc "Redis host"
  configctl add auth.jwt.secret --type string --secret --desc "JWT signing secret"
  configctl validate
  configctl generate-docs`,
		Version: version,
	}

	// Add commands
	rootCmd.AddCommand(newAddCmd())
	rootCmd.AddCommand(newValidateCmd())
	rootCmd.AddCommand(newGenerateDocsCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
