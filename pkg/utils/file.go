package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetModuleDirs returns the root directory of the current module (where go.mod is located)
func GetModuleDirs() ([]string, error) {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "all")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list module directories: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	return lines, nil
}

// GetProjectRoot locates the root directory of the project by traversing up from moduleDir
func GetProjectRoot(moduleDir string, subDir string) (string, error) {
	currentDir := moduleDir
	for {
		// Check if the data directory exists
		if _, err := os.Stat(filepath.Join(currentDir, subDir)); err == nil {
			return currentDir, nil
		}

		// Move up to the parent directory
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// Reached the system root directory (e.g., /), project root not found
			return "", fmt.Errorf("project root not found: %s directory does not exist", subDir)
		}
		currentDir = parentDir
	}
}

func GetFirstExistingFile(subDir string) (string, error) {
	moduleDirs, err := GetModuleDirs()
	if err != nil {
		return "", fmt.Errorf("failed to get module directories: %v", err)
	}
	for _, dir := range moduleDirs {
		projectRoot, err := GetProjectRoot(dir, subDir)
		if err != nil {
			continue // Try the next module directory
		}

		// Check if the file exists in the project root
		filePath := filepath.Join(projectRoot, subDir)
		if _, err := os.Stat(filePath); err == nil {
			return filePath, nil // Return the first existing file path
		}
	}
	return "", fmt.Errorf("project root not found: %s directory does not exist", subDir)
}

// FindAllGoModDirs recursively searches for all directories containing a go.mod file starting from the given rootDir.
// It returns a slice of directory paths where go.mod files are found.
func FindAllGoModDirs(rootDir string) ([]string, error) {
	if _, err := os.Stat(rootDir); err != nil {
		return nil, fmt.Errorf("root directory does not exist: %s", rootDir)
	}

	var goModDirs []string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == "go.mod" {
			dir := filepath.Dir(path)
			goModDirs = append(goModDirs, dir)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return goModDirs, nil
}

// FindAllFiles recursively finds all files (not directories) under the given root directory and returns their paths.
func FindAllFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
