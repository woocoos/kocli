package enum

import (
	"database/sql/driver"
	"fmt"
	"io"
	"strconv"
)

// Code generated by:  kocli gen enum -b int -n namedValuesInt -i -v Name1,2,Name2,3,Name3,4

// NamedValuesInt enum definition
type NamedValuesInt int

const (
	NamedValuesIntName1 NamedValuesInt = 2
	NamedValuesIntName2 NamedValuesInt = 3
	NamedValuesIntName3 NamedValuesInt = 4
)

// String returns the readable string representation of the NamedValuesInt.
func (e NamedValuesInt) String() string {
	switch e {
	case NamedValuesIntName1:
		return "Name1"
	case NamedValuesIntName2:
		return "Name2"
	case NamedValuesIntName3:
		return "Name3"
	default:
		return "Unknown"
	}
}

// Values implements ent EnumValues interface
func (e NamedValuesInt) Values() []string {
	return []string{
		NamedValuesIntName1.String(),
		NamedValuesIntName2.String(),
		NamedValuesIntName3.String(),
	}
}

// MarshalGQL implements graphql.Marshaler.
func (e NamedValuesInt) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(e.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler.
func (e *NamedValuesInt) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enum NamedValuesInt must be strings,but %v", v)
	}
	switch str {
	case "Name1":
		*e = NamedValuesIntName1
	case "Name2":
		*e = NamedValuesIntName2
	case "Name3":
		*e = NamedValuesIntName3
	default:
		return fmt.Errorf("%q is not a valid NamedValuesInt", str)
	}
	return nil
}

// NamedValuesIntValidator validates NamedValuesInt enum value
func NamedValuesIntValidator(et NamedValuesInt) error {
	switch et {
	case NamedValuesIntName1, NamedValuesIntName2, NamedValuesIntName3:
		return nil
	default:
		return fmt.Errorf("invalid enum value for NamedValuesInt field: %q", et)
	}
}

// Value implements driver.Valuer
func (e NamedValuesInt) Value() (driver.Value, error) {
	return int64(e), nil
}

// Scan implements sql.Scanner
func (e *NamedValuesInt) Scan(value any) error {
	switch v := value.(type) {
	case nil:
		return nil
	case int:
		*e = NamedValuesInt(v)
	case int64:
		*e = NamedValuesInt(v)
	default:
		return fmt.Errorf("unsupported type for %s: %T", "NamedValuesInt", value)
	}
	return nil
}
