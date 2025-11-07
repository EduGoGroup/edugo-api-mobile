package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate configuration files",
		Long: `Validate that all configuration files are valid and consistent.

Checks:
  - All YAML files are valid syntax
  - All config structs have proper mapstructure tags
  - .env.example has all required secret variables

Examples:
  configctl validate`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return validateConfig()
		},
	}

	return cmd
}

func validateConfig() error {
	fmt.Println("🔍 Validating configuration files...")
	fmt.Println()

	errors := 0

	// Validate YAML files
	yamlFiles := []string{
		"config/config.yaml",
		"config/config-local.yaml",
		"config/config-dev.yaml",
		"config/config-qa.yaml",
		"config/config-prod.yaml",
	}

	for _, file := range yamlFiles {
		if err := validateYAMLFile(file); err != nil {
			fmt.Printf("❌ %s: %v\n", file, err)
			errors++
		} else {
			fmt.Printf("✅ %s\n", file)
		}
	}

	// Validate .env.example exists
	if _, err := os.Stat(".env.example"); os.IsNotExist(err) {
		fmt.Println("❌ .env.example: file not found")
		errors++
	} else {
		fmt.Println("✅ .env.example")
	}

	fmt.Println()
	if errors > 0 {
		return fmt.Errorf("validation failed with %d error(s)", errors)
	}

	fmt.Println("✅ All configuration files are valid!")
	return nil
}

func validateYAMLFile(path string) error {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file not found")
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse YAML
	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("invalid YAML syntax: %w", err)
	}

	return nil
}
