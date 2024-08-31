package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/woocoos/knockout-go/ent/schemax"
)

type Resource struct {
	ent.Schema
}

func (Resource) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "test_resource"},
		schemax.Resources([]string{"name"}),
		schemax.TenantField("tenant_id"),
	}
}

func (Resource) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("description").Optional(),
		field.String("tenant_id").Optional(),
	}
}
