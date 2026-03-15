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

func WriteBarrelFiles(categories []BarrelFileCategory) error {
	for _, cat := range categories {
		if len(cat.DefinitionNames) == 0 {
			continue
		}
		if err := writeBarrelFile(cat.DirPath, cat.DefinitionNames); err != nil {
			return fmt.Errorf("failed to write barrel file in %s: %w", cat.DirPath, err)
		}
	}
	return nil
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
