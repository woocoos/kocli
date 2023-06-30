// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/errcode"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/woocoos/entco/pkg/pagination"
	"github.com/woocoos/kocli/integration/resource/ent/resource"
)

// Common entgql types.
type (
	Cursor         = entgql.Cursor[int]
	PageInfo       = entgql.PageInfo[int]
	OrderDirection = entgql.OrderDirection
)

func orderFunc(o OrderDirection, field string) func(*sql.Selector) {
	if o == entgql.OrderDirectionDesc {
		return Desc(field)
	}
	return Asc(field)
}

const errInvalidPagination = "INVALID_PAGINATION"

func validateFirstLast(first, last *int) (err *gqlerror.Error) {
	switch {
	case first != nil && last != nil:
		err = &gqlerror.Error{
			Message: "Passing both `first` and `last` to paginate a connection is not supported.",
		}
	case first != nil && *first < 0:
		err = &gqlerror.Error{
			Message: "`first` on a connection cannot be less than zero.",
		}
		errcode.Set(err, errInvalidPagination)
	case last != nil && *last < 0:
		err = &gqlerror.Error{
			Message: "`last` on a connection cannot be less than zero.",
		}
		errcode.Set(err, errInvalidPagination)
	}
	return err
}

func collectedField(ctx context.Context, path ...string) *graphql.CollectedField {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return nil
	}
	field := fc.Field
	oc := graphql.GetOperationContext(ctx)
walk:
	for _, name := range path {
		for _, f := range graphql.CollectFields(oc, field.Selections, nil) {
			if f.Alias == name {
				field = f
				continue walk
			}
		}
		return nil
	}
	return &field
}

func hasCollectedField(ctx context.Context, path ...string) bool {
	if graphql.GetFieldContext(ctx) == nil {
		return true
	}
	return collectedField(ctx, path...) != nil
}

const (
	edgesField      = "edges"
	nodeField       = "node"
	pageInfoField   = "pageInfo"
	totalCountField = "totalCount"
)

func paginateLimit(first, last *int) int {
	var limit int
	if first != nil {
		limit = *first + 1
	} else if last != nil {
		limit = *last + 1
	}
	return limit
}

// ResourceEdge is the edge representation of Resource.
type ResourceEdge struct {
	Node   *Resource `json:"node"`
	Cursor Cursor    `json:"cursor"`
}

// ResourceConnection is the connection containing edges to Resource.
type ResourceConnection struct {
	Edges      []*ResourceEdge `json:"edges"`
	PageInfo   PageInfo        `json:"pageInfo"`
	TotalCount int             `json:"totalCount"`
}

func (c *ResourceConnection) build(nodes []*Resource, pager *resourcePager, after *Cursor, first *int, before *Cursor, last *int) {
	c.PageInfo.HasNextPage = before != nil
	c.PageInfo.HasPreviousPage = after != nil
	if first != nil && *first+1 == len(nodes) {
		c.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && *last+1 == len(nodes) {
		c.PageInfo.HasPreviousPage = true
		nodes = nodes[:len(nodes)-1]
	}
	var nodeAt func(int) *Resource
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *Resource {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *Resource {
			return nodes[i]
		}
	}
	c.Edges = make([]*ResourceEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		c.Edges[i] = &ResourceEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}
	if l := len(c.Edges); l > 0 {
		c.PageInfo.StartCursor = &c.Edges[0].Cursor
		c.PageInfo.EndCursor = &c.Edges[l-1].Cursor
	}
	if c.TotalCount == 0 {
		c.TotalCount = len(nodes)
	}
}

// ResourcePaginateOption enables pagination customization.
type ResourcePaginateOption func(*resourcePager) error

// WithResourceOrder configures pagination ordering.
func WithResourceOrder(order *ResourceOrder) ResourcePaginateOption {
	if order == nil {
		order = DefaultResourceOrder
	}
	o := *order
	return func(pager *resourcePager) error {
		if err := o.Direction.Validate(); err != nil {
			return err
		}
		if o.Field == nil {
			o.Field = DefaultResourceOrder.Field
		}
		pager.order = &o
		return nil
	}
}

// WithResourceFilter configures pagination filter.
func WithResourceFilter(filter func(*ResourceQuery) (*ResourceQuery, error)) ResourcePaginateOption {
	return func(pager *resourcePager) error {
		if filter == nil {
			return errors.New("ResourceQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type resourcePager struct {
	reverse bool
	order   *ResourceOrder
	filter  func(*ResourceQuery) (*ResourceQuery, error)
}

func newResourcePager(opts []ResourcePaginateOption, reverse bool) (*resourcePager, error) {
	pager := &resourcePager{reverse: reverse}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultResourceOrder
	}
	return pager, nil
}

func (p *resourcePager) applyFilter(query *ResourceQuery) (*ResourceQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *resourcePager) toCursor(r *Resource) Cursor {
	return p.order.Field.toCursor(r)
}

func (p *resourcePager) applyCursors(query *ResourceQuery, after, before *Cursor) (*ResourceQuery, error) {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	for _, predicate := range entgql.CursorsPredicate(after, before, DefaultResourceOrder.Field.column, p.order.Field.column, direction) {
		query = query.Where(predicate)
	}
	return query, nil
}

func (p *resourcePager) applyOrder(query *ResourceQuery) *ResourceQuery {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	query = query.Order(p.order.Field.toTerm(direction.OrderTermOption()))
	if p.order.Field != DefaultResourceOrder.Field {
		query = query.Order(DefaultResourceOrder.Field.toTerm(direction.OrderTermOption()))
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(p.order.Field.column)
	}
	return query
}

func (p *resourcePager) orderExpr(query *ResourceQuery) sql.Querier {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(p.order.Field.column)
	}
	return sql.ExprFunc(func(b *sql.Builder) {
		b.Ident(p.order.Field.column).Pad().WriteString(string(direction))
		if p.order.Field != DefaultResourceOrder.Field {
			b.Comma().Ident(DefaultResourceOrder.Field.column).Pad().WriteString(string(direction))
		}
	})
}

// Paginate executes the query and returns a relay based cursor connection to Resource.
func (r *ResourceQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...ResourcePaginateOption,
) (*ResourceConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newResourcePager(opts, last != nil)
	if err != nil {
		return nil, err
	}
	if r, err = pager.applyFilter(r); err != nil {
		return nil, err
	}
	conn := &ResourceConnection{Edges: []*ResourceEdge{}}
	ignoredEdges := !hasCollectedField(ctx, edgesField)
	if hasCollectedField(ctx, totalCountField) || hasCollectedField(ctx, pageInfoField) {
		hasPagination := after != nil || first != nil || before != nil || last != nil
		if hasPagination || ignoredEdges {
			if conn.TotalCount, err = r.Clone().Count(ctx); err != nil {
				return nil, err
			}
			conn.PageInfo.HasNextPage = first != nil && conn.TotalCount > 0
			conn.PageInfo.HasPreviousPage = last != nil && conn.TotalCount > 0
		}
	}
	if ignoredEdges || (first != nil && *first == 0) || (last != nil && *last == 0) {
		return conn, nil
	}
	if r, err = pager.applyCursors(r, after, before); err != nil {
		return nil, err
	}
	if limit := paginateLimit(first, last); limit != 0 {
		r.Limit(limit)
	}
	if sp, ok := pagination.SimplePaginationFromContext(ctx); ok {
		if first != nil {
			r.Offset((sp.PageIndex - sp.CurrentIndex - 1) * *first)
		}
		if last != nil {
			r.Offset((sp.CurrentIndex - sp.PageIndex - 1) * *last)
		}
	}
	if field := collectedField(ctx, edgesField, nodeField); field != nil {
		if err := r.collectField(ctx, graphql.GetOperationContext(ctx), *field, []string{edgesField, nodeField}); err != nil {
			return nil, err
		}
	}
	r = pager.applyOrder(r)
	nodes, err := r.All(ctx)
	if err != nil {
		return nil, err
	}
	conn.build(nodes, pager, after, first, before, last)
	return conn, nil
}

// ResourceOrderField defines the ordering field of Resource.
type ResourceOrderField struct {
	// Value extracts the ordering value from the given Resource.
	Value    func(*Resource) (ent.Value, error)
	column   string // field or computed.
	toTerm   func(...sql.OrderTermOption) resource.OrderOption
	toCursor func(*Resource) Cursor
}

// ResourceOrder defines the ordering of Resource.
type ResourceOrder struct {
	Direction OrderDirection      `json:"direction"`
	Field     *ResourceOrderField `json:"field"`
}

// DefaultResourceOrder is the default ordering of Resource.
var DefaultResourceOrder = &ResourceOrder{
	Direction: entgql.OrderDirectionAsc,
	Field: &ResourceOrderField{
		Value: func(r *Resource) (ent.Value, error) {
			return r.ID, nil
		},
		column: resource.FieldID,
		toTerm: resource.ByID,
		toCursor: func(r *Resource) Cursor {
			return Cursor{ID: r.ID}
		},
	},
}

// ToEdge converts Resource into ResourceEdge.
func (r *Resource) ToEdge(order *ResourceOrder) *ResourceEdge {
	if order == nil {
		order = DefaultResourceOrder
	}
	return &ResourceEdge{
		Node:   r,
		Cursor: order.Field.toCursor(r),
	}
}