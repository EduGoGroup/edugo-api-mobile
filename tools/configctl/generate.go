package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newGenerateDocsCmd() *cobra.Command {
	var outputFile string

	cmd := &cobra.Command{
		Use:   "generate-docs",
		Short: "Generate configuration documentation",
		Long: `Generate CONFIG.md with documentation of all configuration variables.

The documentation includes:
  - Variable name and type
  - Description
  - Default value
  - Whether it's required
  - Environment variable mapping

Examples:
  configctl generate-docs
  configctl generate-docs --output docs/CONFIG.md`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateDocs(outputFile)
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "CONFIG.md", "Output file path")

	return cmd
}

func generateDocs(outputFile string) error {
	fmt.Printf("📚 Generating documentation to %s...\n", outputFile)
	fmt.Println()

	// TODO: Implement actual logic
	fmt.Println("⚠️  This command is not yet fully implemented")
	fmt.Println("   The following would be generated:")
	fmt.Println("   - Table of all configuration variables")
	fmt.Println("   - ENV var → config path mapping")
	fmt.Println("   - Examples for each environment")

	return nil
}
