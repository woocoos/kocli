package enum

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAfterGenerate(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		var iv IntBase
		err := iv.UnmarshalGQL("value1")
		assert.NoError(t, err)
		assert.Equal(t, IntBaseValue1, iv)
		err = iv.UnmarshalGQL("Value2")
		assert.NoError(t, err)
		assert.Equal(t, IntBaseValue2, iv)
		err = iv.UnmarshalGQL("Value3")
		assert.NoError(t, err)
		assert.Equal(t, IntBaseValue3, iv)
		err = iv.UnmarshalGQL("Value4")
		assert.ErrorContains(t, err, "not a valid IntBase", err.Error())
		assert.Equal(t, "Value3", iv.String())
		assert.Equal(t, []string{"value1", "Value2", "Value3"}, iv.Values())
	})
	t.Run("string", func(t *testing.T) {
		var sv StringBase
		err := sv.UnmarshalGQL("value1")
		assert.NoError(t, err)
		assert.Equal(t, StringBaseValue1, sv)
		err = sv.UnmarshalGQL("Value2")
		assert.NoError(t, err)
		assert.Equal(t, StringBaseValue2, sv)
		err = sv.UnmarshalGQL("Value4")
		assert.ErrorContains(t, err, "not a valid StringBase", err.Error())
		assert.Equal(t, "Value2", sv.String())
		assert.Equal(t, []string{"value1", "Value2", "Value3"}, sv.Values())
	})
	t.Run("named_values", func(t *testing.T) {
		var nv NamedValues
		err := nv.UnmarshalGQL("Name1")
		assert.NoError(t, err)
		assert.Equal(t, NamedValuesName1, nv)
		err = nv.UnmarshalGQL("Name2")
		assert.NoError(t, err)
		assert.Equal(t, NamedValuesName2, nv)
		err = nv.UnmarshalGQL("Name4")
		assert.ErrorContains(t, err, "not a valid NamedValues", err.Error())
		assert.Equal(t, "Name2", nv.String())
		assert.Equal(t, []string{"Name1", "Name2", "Name3"}, nv.Values())
	})
	t.Run("named_values_int", func(t *testing.T) {
		var nv NamedValuesInt
		err := nv.UnmarshalGQL("Name1")
		assert.NoError(t, err)
		assert.Equal(t, NamedValuesIntName1, nv)
		err = nv.UnmarshalGQL("Name2")
		assert.NoError(t, err)
		assert.Equal(t, NamedValuesIntName2, nv)
		err = nv.UnmarshalGQL("Name4")
		assert.ErrorContains(t, err, "not a valid NamedValuesInt", err.Error())
		assert.Equal(t, "Name2", nv.String())
		assert.Equal(t, []string{"Name1", "Name2", "Name3"}, nv.Values())
	})
}
