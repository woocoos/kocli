// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// TestResourceColumns holds the columns for the "test_resource" table.
	TestResourceColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "tenant_id", Type: field.TypeString, Nullable: true},
	}
	// TestResourceTable holds the schema information for the "test_resource" table.
	TestResourceTable = &schema.Table{
		Name:       "test_resource",
		Columns:    TestResourceColumns,
		PrimaryKey: []*schema.Column{TestResourceColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		TestResourceTable,
	}
)

func init() {
	TestResourceTable.Annotation = &entsql.Annotation{
		Table: "test_resource",
	}
}
