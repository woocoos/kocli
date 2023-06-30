// Code generated by ent, DO NOT EDIT.

package ent

import (
	"github.com/woocoos/kocli/integration/resource/ent/resource"
	"github.com/woocoos/kocli/integration/resource/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	resourceFields := schema.Resource{}.Fields()
	_ = resourceFields
	// resourceDescName is the schema descriptor for name field.
	resourceDescName := resourceFields[0].Descriptor()
	// resource.NameValidator is a validator for the "name" field. It is called by the builders before save.
	resource.NameValidator = resourceDescName.Validators[0].(func(string) error)
}