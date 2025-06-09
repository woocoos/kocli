package gentool

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tsingsun/woocoo/cmd/woco/gen"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// EnumInput represents the input for generating an enum.
type EnumInput struct {
	TargetDir     string
	BaseType      string
	EnumName      string
	InputValues   []string
	IsNamedValues bool
	Names         []string // 添加Names字段
	Values        []string // 添加ValuesSlice字段
}

// GenerateEnum generates the enum code based on the provided input.
func GenerateEnum(input EnumInput) error {
	if err := enumInput(&input); err != nil {
		return err
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
		outputPath := filepath.Join(input.TargetDir, fmt.Sprintf("%s.go", sf(input.EnumName)))
		if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
			return err
		}
		if err := gen.FormatGoFile(outputPath, b.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

func enumInput(input *EnumInput) (err error) {
	if input.TargetDir == "" {
		input.TargetDir = filepath.Join("codegen", "entgen", "types")
	}
	input.TargetDir, err = filepath.Abs(input.TargetDir)
	if err != nil {
		return err
	}
	input.Names = make([]string, 0)
	input.Values = make([]string, 0)
	if input.IsNamedValues {
		tv := input.InputValues
		if len(tv)%2 != 0 {
			return errors.New("values length must be even")
		}
		for i := 0; i < len(tv); i += 2 {
			input.Names = append(input.Names, tv[i])
			switch input.BaseType {
			case "int":
				input.Values = append(input.Values, tv[i+1])
			default:
				input.Values = append(input.Values, strconv.Quote(tv[i+1]))
			}
		}
	} else {
		// 转化为NamedValues
		for i, s := range input.InputValues {
			input.Names = append(input.Names, s)
			switch input.BaseType {
			case "int":
				input.Values = append(input.Values, strconv.Itoa(i+1))
			default:
				input.Values = append(input.Values, strconv.Quote(s))
			}
		}
	}
	return nil
}
