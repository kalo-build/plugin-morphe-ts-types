package compile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kalo-build/go-util/strcase"
)

type BarrelFileCategory struct {
	DirPath         string
	DefinitionNames []string
}

func WriteBarrelFiles(rootDirPath string, categories []BarrelFileCategory) error {
	for _, cat := range categories {
		if len(cat.DefinitionNames) == 0 {
			continue
		}
		if err := writeBarrelFile(cat.DirPath, cat.DefinitionNames); err != nil {
			return fmt.Errorf("failed to write barrel file in %s: %w", cat.DirPath, err)
		}
	}
	if rootDirPath != "" {
		if err := writeRootBarrelFile(rootDirPath, categories); err != nil {
			return fmt.Errorf("failed to write root barrel file: %w", err)
		}
	}
	return nil
}

func writeRootBarrelFile(rootDirPath string, categories []BarrelFileCategory) error {
	var categoryDirs []string
	for _, cat := range categories {
		if len(cat.DefinitionNames) == 0 {
			continue
		}
		relDir, err := filepath.Rel(rootDirPath, cat.DirPath)
		if err != nil {
			continue
		}
		categoryDirs = append(categoryDirs, filepath.ToSlash(relDir))
	}
	if len(categoryDirs) == 0 {
		return nil
	}
	sort.Strings(categoryDirs)

	var lines []string
	for _, dir := range categoryDirs {
		lines = append(lines, fmt.Sprintf("export * from './%s'", dir))
	}

	content := strings.Join(lines, "\n") + "\n"
	barrelPath := filepath.Join(rootDirPath, "index.d.ts")

	if err := os.MkdirAll(rootDirPath, 0755); err != nil {
		return err
	}
	return os.WriteFile(barrelPath, []byte(content), 0644)
}

func writeBarrelFile(dirPath string, definitionNames []string) error {
	sorted := make([]string, len(definitionNames))
	copy(sorted, definitionNames)
	sort.Strings(sorted)

	var lines []string
	for _, name := range sorted {
		kebab := strcase.ToKebabCaseLower(name)
		lines = append(lines, fmt.Sprintf("export * from './%s'", kebab))
	}

	content := strings.Join(lines, "\n") + "\n"
	barrelPath := filepath.Join(dirPath, "index.d.ts")

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}
	return os.WriteFile(barrelPath, []byte(content), 0644)
}
