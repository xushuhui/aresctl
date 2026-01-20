package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aresctl",
	Short: "Aresctl - Command-line tool for Ares framework",
	Long:  `Aresctl is a powerful CLI tool that helps you develop applications with the Ares framework. It provides code generation, OpenAPI documentation, and more.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(openapiCmd)
}

var openapiCmd = &cobra.Command{
	Use:   "openapi",
	Short: "Generate OpenAPI specification",
	Long:  `Generate an openapi.yaml file from your Go code by analyzing route definitions and API structures`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := generateOpenAPI(); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating OpenAPI: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Generated openapi.yaml successfully")
	},
}

func generateOpenAPI() error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Default paths - can be made configurable via flags
	routeDir := cwd + "/internal/server"
	apiDir := cwd + "/api"
	outputFile := cwd + "/openapi.yaml"

	// Check if directories exist
	if _, err := os.Stat(routeDir); os.IsNotExist(err) {
		return fmt.Errorf("route directory not found: %s", routeDir)
	}
	if _, err := os.Stat(apiDir); os.IsNotExist(err) {
		return fmt.Errorf("api directory not found: %s", apiDir)
	}

	// Generate OpenAPI spec
	GenerateOpenAPI(routeDir, apiDir, outputFile)

	return nil
}
