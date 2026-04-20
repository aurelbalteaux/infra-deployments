// Command validate-refs checks that all YAML files in a directory tree are
// referenced in their parent kustomization.yaml files.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/redhat-appstudio/infra-deployments/infra-tools/internal/kustomize"
)

// version is set via -ldflags at build time.
var version = "dev"

func main() {
	var (
		rootDir     = flag.String("root", "", "Root directory to validate (required)")
		showVersion = flag.Bool("version", false, "Print version and exit")
		verbose     = flag.Bool("verbose", false, "Show all checked directories")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("validate-refs %s\n", version)
		os.Exit(0)
	}

	if *rootDir == "" {
		fmt.Fprintf(os.Stderr, "Error: --root is required\n")
		flag.Usage()
		os.Exit(1)
	}

	absRoot, err := filepath.Abs(*rootDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving root directory: %v\n", err)
		os.Exit(1)
	}

	// Check that the root directory exists
	if info, err := os.Stat(absRoot); err != nil || !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: %s is not a valid directory\n", absRoot)
		os.Exit(1)
	}

	fmt.Printf("Validating YAML references in: %s\n\n", absRoot)

	result, err := kustomize.ValidateAllReferences(absRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during validation: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Checked %d directories with kustomization files\n\n", result.CheckedDirs)
	}

	if len(result.OrphanedFiles) == 0 {
		fmt.Println("✓ All YAML files are properly referenced in kustomization files")
		return
	}

	// Group orphaned files by directory for better readability
	byDir := make(map[string][]kustomize.OrphanedFile)
	for _, orphan := range result.OrphanedFiles {
		byDir[orphan.KustomizeDir] = append(byDir[orphan.KustomizeDir], orphan)
	}

	fmt.Printf("✗ Found %d orphaned YAML file(s):\n\n", len(result.OrphanedFiles))

	for dir, orphans := range byDir {
		relDir, _ := filepath.Rel(absRoot, dir)
		if relDir == "" {
			relDir = "."
		}
		fmt.Printf("  Directory: %s/\n", relDir)
		for _, orphan := range orphans {
			fmt.Printf("    - %s\n", filepath.Base(orphan.Path))
		}
		fmt.Println()
	}

	fmt.Printf("These files should be added to their respective kustomization.yaml files\n")
	fmt.Printf("or removed if they are no longer needed.\n")

	os.Exit(1)
}
