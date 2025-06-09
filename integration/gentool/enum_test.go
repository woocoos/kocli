package gentool

import (
	"github.com/woocoos/kocli/gentool"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGenerateEnum tests the GenerateEnum function.
func TestGenerateEnum(t *testing.T) {
	t.Run("int base", func(t *testing.T) {
		// Define test input
		input := gentool.EnumInput{
			TargetDir:   "./enum",
			BaseType:    "int",
			EnumName:    "intBase",
			InputValues: []string{"value1", "Value2", "Value3"},
		}
		err := gentool.GenerateEnum(input)
		assert.NoError(t, err)
	})
	t.Run("string base", func(t *testing.T) {
		input := gentool.EnumInput{
			TargetDir:   "./enum",
			BaseType:    "string",
			EnumName:    "stringBase",
			InputValues: []string{"value1", "Value2", "Value3"},
		}
		err := gentool.GenerateEnum(input)
		assert.NoError(t, err)
	})
	t.Run("namedValues", func(t *testing.T) {
		input := gentool.EnumInput{
			TargetDir:     "./enum",
			BaseType:      "string",
			EnumName:      "namedValues",
			InputValues:   []string{"Name1", "value1", "Name2", "Value2", "Name3", "Value3"},
			IsNamedValues: true,
		}
		err := gentool.GenerateEnum(input)
		assert.NoError(t, err)
	})
	t.Run("namedValuesInt", func(t *testing.T) {
		input := gentool.EnumInput{
			TargetDir:     "./enum",
			BaseType:      "int",
			EnumName:      "namedValuesInt",
			InputValues:   []string{"Name1", "2", "Name2", "3", "Name3", "4"},
			IsNamedValues: true,
		}
		err := gentool.GenerateEnum(input)
		assert.NoError(t, err)
	})
}
