package main

import (
	"fmt"
	"os"

	"github.com/goharbor/harbor-helm-migration-tool/pkg/migrate"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	dryRun  bool
	verbose bool
	strict  bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "harbor-migrate <input-values.yaml> [output-values.yaml]",
		Short: "Migrate Harbor Helm chart values from harbor-helm to harbor-next-helm format",
		Long: `harbor-migrate transforms values.yaml files from the legacy harbor-helm chart
to the new harbor-next-helm format.

It handles field renames, structural changes, and generates warnings for
features that require manual migration or are no longer supported.

Examples:
  # Migrate and write to new file
  harbor-migrate legacy-values.yaml new-values.yaml

  # Preview migration without writing (dry-run)
  harbor-migrate --dry-run legacy-values.yaml

  # Migrate with verbose output
  harbor-migrate --verbose legacy-values.yaml new-values.yaml

  # Strict mode - exit with error if warnings found
  harbor-migrate --strict legacy-values.yaml new-values.yaml`,
		Args:    cobra.RangeArgs(1, 2),
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
		RunE:    runMigrate,
	}

	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print migrated values to stdout without writing file")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().BoolVar(&strict, "strict", false, "Exit with error if any warnings are generated")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}

func runMigrate(cmd *cobra.Command, args []string) error {
	inputFile := args[0]
	var outputFile string
	if len(args) > 1 {
		outputFile = args[1]
	}

	// Read input file
	if verbose {
		fmt.Fprintf(os.Stderr, "Reading input file: %s\n", inputFile)
	}

	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Perform migration
	if verbose {
		fmt.Fprintln(os.Stderr, "Performing migration...")
	}

	result, err := migrate.MigrateFromYAML(data)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	// Print warnings to stderr
	migrate.PrintWarnings(os.Stderr, result.Warnings)

	// Generate output YAML
	var output []byte
	if dryRun || outputFile == "" {
		output, err = result.ToYAML()
	} else {
		output, err = result.ToYAMLWithHeader(inputFile)
	}
	if err != nil {
		return fmt.Errorf("failed to generate output YAML: %w", err)
	}

	// Output handling
	if dryRun || outputFile == "" {
		// Print to stdout
		fmt.Print(string(output))
	} else {
		// Write to file
		if verbose {
			fmt.Fprintf(os.Stderr, "Writing output file: %s\n", outputFile)
		}

		if err := os.WriteFile(outputFile, output, 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}

		fmt.Fprintf(os.Stderr, "\nMigrated values written to: %s\n", outputFile)
	}

	// Handle strict mode
	if strict && migrate.HasWarnings(result.Warnings) {
		return fmt.Errorf("migration completed with warnings (strict mode enabled)")
	}

	// Exit code 1 if there are errors (but not in strict mode for warnings)
	if migrate.HasErrors(result.Warnings) {
		os.Exit(1)
	}

	return nil
}
