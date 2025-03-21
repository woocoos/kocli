package gentool

import (
	"bytes"
	"fmt"
	"github.com/tsingsun/woocoo/cmd/woco/gen"
	"os"
	"path/filepath"
	"strings"
)

// EnumInput represents the input for generating an enum.
type EnumInput struct {
	targetDir string
	BaseType  string
	EnumName  string
	Values    []string
}

// GenerateEnum generates the enum code based on the provided input.
func GenerateEnum(input EnumInput) error {
	if input.targetDir == "" {
		input.targetDir = filepath.Join("codegen", "entgen", "types")
	}
	// Define the template
	tmpls := gen.ParseT("template/enum.tmpl", templateDir, gen.Funcs, "template/enum.tmpl")
	for _, tpl := range tmpls.Templates() {
		if strings.HasSuffix(tpl.Name(), ".tmpl") {
			// skip
			continue
		}

		b := new(bytes.Buffer)
		// Execute the template
		if err := tpl.Execute(b, input); err != nil {
			return err
		}
		sf := gen.Funcs["snake"].(func(string) string)
		outputPath := filepath.Join(input.targetDir, fmt.Sprintf("%s.go", sf(input.EnumName)))
		if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
			return err
		}
		if err := gen.FormatGoFile(outputPath, b.Bytes()); err != nil {
			return err
		}
	}
	return nil
}
