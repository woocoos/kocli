// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/woocoos/entco/pkg/pagination"
	"github.com/woocoos/kocli/integration/resource/ent/resource"
)

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (r *ResourceQuery) CollectFields(ctx context.Context, satisfies ...string) (*ResourceQuery, error) {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return r, nil
	}
	if err := r.collectField(ctx, graphql.GetOperationContext(ctx), fc.Field, nil, satisfies...); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *ResourceQuery) collectField(ctx context.Context, opCtx *graphql.OperationContext, collected graphql.CollectedField, path []string, satisfies ...string) error {
	path = append([]string(nil), path...)
	var (
		unknownSeen    bool
		fieldSeen      = make(map[string]struct{}, len(resource.Columns))
		selectedFields = []string{resource.FieldID}
	)
	for _, field := range graphql.CollectFields(opCtx, collected.Selections, satisfies) {
		switch field.Name {
		case "name":
			if _, ok := fieldSeen[resource.FieldName]; !ok {
				selectedFields = append(selectedFields, resource.FieldName)
				fieldSeen[resource.FieldName] = struct{}{}
			}
		case "description":
			if _, ok := fieldSeen[resource.FieldDescription]; !ok {
				selectedFields = append(selectedFields, resource.FieldDescription)
				fieldSeen[resource.FieldDescription] = struct{}{}
			}
		case "tenantID":
			if _, ok := fieldSeen[resource.FieldTenantID]; !ok {
				selectedFields = append(selectedFields, resource.FieldTenantID)
				fieldSeen[resource.FieldTenantID] = struct{}{}
			}
		case "id":
		case "__typename":
		default:
			unknownSeen = true
		}
	}
	if !unknownSeen {
		r.Select(selectedFields...)
	}
	return nil
}

type resourcePaginateArgs struct {
	first, last   *int
	after, before *Cursor
	opts          []ResourcePaginateOption
}

func newResourcePaginateArgs(rv map[string]any) *resourcePaginateArgs {
	args := &resourcePaginateArgs{}
	if rv == nil {
		return args
	}
	if v := rv[firstField]; v != nil {
		args.first = v.(*int)
	}
	if v := rv[lastField]; v != nil {
		args.last = v.(*int)
	}
	if v := rv[afterField]; v != nil {
		args.after = v.(*Cursor)
	}
	if v := rv[beforeField]; v != nil {
		args.before = v.(*Cursor)
	}
	if v, ok := rv[whereField].(*ResourceWhereInput); ok {
		args.opts = append(args.opts, WithResourceFilter(v.Filter))
	}
	return args
}

const (
	afterField     = "after"
	firstField     = "first"
	beforeField    = "before"
	lastField      = "last"
	orderByField   = "orderBy"
	directionField = "direction"
	fieldField     = "field"
	whereField     = "where"
)

func fieldArgs(ctx context.Context, whereInput any, path ...string) map[string]any {
	field := collectedField(ctx, path...)
	if field == nil || field.Arguments == nil {
		return nil
	}
	oc := graphql.GetOperationContext(ctx)
	args := field.ArgumentMap(oc.Variables)
	return unmarshalArgs(ctx, whereInput, args)
}

// unmarshalArgs allows extracting the field arguments from their raw representation.
func unmarshalArgs(ctx context.Context, whereInput any, args map[string]any) map[string]any {
	for _, k := range []string{firstField, lastField} {
		v, ok := args[k]
		if !ok {
			continue
		}
		i, err := graphql.UnmarshalInt(v)
		if err == nil {
			args[k] = &i
		}
	}
	for _, k := range []string{beforeField, afterField} {
		v, ok := args[k]
		if !ok {
			continue
		}
		c := &Cursor{}
		if c.UnmarshalGQL(v) == nil {
			args[k] = c
		}
	}
	if v, ok := args[whereField]; ok && whereInput != nil {
		if err := graphql.UnmarshalInputFromContext(ctx, v, whereInput); err == nil {
			args[whereField] = whereInput
		}
	}

	return args
}

func limitRows(ctx context.Context, partitionBy string, limit int, first, last *int, orderBy ...sql.Querier) func(s *sql.Selector) {
	offset := 0
	if sp, ok := pagination.SimplePaginationFromContext(ctx); ok {
		if first != nil {
			offset = (sp.PageIndex - sp.CurrentIndex - 1) * *first
		}
		if last != nil {
			offset = (sp.CurrentIndex - sp.PageIndex - 1) * *last
		}
	}
	return func(s *sql.Selector) {
		d := sql.Dialect(s.Dialect())
		s.SetDistinct(false)
		with := d.With("src_query").
			As(s.Clone()).
			With("limited_query").
			As(
				d.Select("*").
					AppendSelectExprAs(
						sql.RowNumber().PartitionBy(partitionBy).OrderExpr(orderBy...),
						"row_number",
					).
					From(d.Table("src_query")),
			)
		t := d.Table("limited_query").As(s.TableName())
		if offset != 0 {
			*s = *d.Select(s.UnqualifiedColumns()...).
				From(t).
				Where(sql.GT(t.C("row_number"), offset)).Limit(limit).
				Prefix(with)
		} else {
			*s = *d.Select(s.UnqualifiedColumns()...).
				From(t).
				Where(sql.LTE(t.C("row_number"), limit)).
				Prefix(with)
		}
	}
}

// mayAddCondition appends another type condition to the satisfies list
// if condition is enabled (Node/Nodes) and it does not exist in the list.
func mayAddCondition(satisfies []string, typeCond string) []string {
	if len(satisfies) == 0 {
		return satisfies
	}
	for _, s := range satisfies {
		if typeCond == s {
			return satisfies
		}
	}
	return append(satisfies, typeCond)
}
