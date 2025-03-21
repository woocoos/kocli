package gentool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGenerateEnum tests the GenerateEnum function.
func TestGenerateEnum(t *testing.T) {
	t.Run("int base", func(t *testing.T) {
		// Define test input
		input := EnumInput{
			targetDir: "testdata/tmp",
			BaseType:  "int",
			EnumName:  "intBase",
			Values:    []string{"Value1", "Value2", "Value3"},
		}

		// Call the function
		err := GenerateEnum(input)
		assert.NoError(t, err)
	})
	t.Run("string base", func(t *testing.T) {
		input := EnumInput{
			targetDir: "testdata/tmp",
			BaseType:  "string",
			EnumName:  "stringBase",
			Values:    []string{"Value1", "Value2", "Value3"},
		}
		// Call the function
		err := GenerateEnum(input)
		assert.NoError(t, err)
	})
}
