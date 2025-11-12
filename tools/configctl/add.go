package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {
	var (
		varType     string
		isSecret    bool
		defaultVal  string
		description string
		dryRun      bool
	)

	cmd := &cobra.Command{
		Use:   "add [hierarchy.path]",
		Short: "Add a new configuration variable",
		Long: `Add a new configuration variable to the system.

The tool will:
  - Update the Go config struct in internal/config/config.go
  - Update YAML files (if public) or .env.example (if secret)
  - Update validator.go (if secret and required)

Examples:
  # Add a public config variable
  configctl add database.postgres.pool_size --type int --default 10 --desc "Connection pool size"

  # Add a secret variable
  configctl add auth.jwt.secret --type string --secret --desc "JWT signing secret"

  # Dry run to preview changes
  configctl add storage.s3.timeout --type duration --default 30s --desc "S3 timeout" --dry-run`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			if dryRun {
				fmt.Println("🔍 DRY RUN MODE - No changes will be made")
				fmt.Println()
			}

			return addVariable(path, varType, isSecret, defaultVal, description, dryRun)
		},
	}

	cmd.Flags().StringVar(&varType, "type", "string", "Variable type (string, int, bool, duration)")
	cmd.Flags().BoolVar(&isSecret, "secret", false, "Mark as secret (ENV only)")
	cmd.Flags().StringVar(&defaultVal, "default", "", "Default value (for public vars)")
	cmd.Flags().StringVar(&description, "desc", "", "Description of the variable")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without applying them")
	_ = cmd.MarkFlagRequired("desc")

	return cmd
}

func addVariable(path, varType string, isSecret bool, defaultVal, description string, dryRun bool) error {
	fmt.Printf("📝 Adding variable: %s\n", path)
	fmt.Printf("   Type: %s\n", varType)
	fmt.Printf("   Secret: %v\n", isSecret)
	if defaultVal != "" {
		fmt.Printf("   Default: %s\n", defaultVal)
	}
	fmt.Printf("   Description: %s\n", description)
	fmt.Println()

	// TODO: Implement actual logic
	fmt.Println("⚠️  This command is not yet fully implemented")
	fmt.Println("   The following steps would be performed:")
	fmt.Println("   1. Validate hierarchy path")
	fmt.Println("   2. Update internal/config/config.go")
	if isSecret {
		fmt.Println("   3. Update .env.example")
		fmt.Println("   4. Update internal/config/validator.go")
	} else {
		fmt.Println("   3. Update config/*.yaml files")
	}

	return nil
}
