package engine

// SortDirection represents a query sort SortDirection
type SortDirection byte

const (
	// Ascending means going up, A-Z
	Ascending SortDirection = 1 << iota

	// Descending means reverse order, Z-A
	Descending
)

// CompareCondition represents a filter comparison operation
// between a field and a value
type CompareCondition byte

const (
	// Equal if it should be the same
	Equal CompareCondition = 1 << iota

	// LessThan if it should be smaller
	LessThan

	// LessThanOrEqual if it should be smaller or equal
	LessThanOrEqual

	// GreaterThan if it should be larger
	GreaterThan

	// GreaterThanOrEqual if it should be equal or greater than
	GreaterThanOrEqual
)

type (
	// Query represents a query specification for filtering
	// sorting, paging and limiting the data requested
	Query struct {
		Name    string
		Offset  int
		Limit   int
		Filters []*FilterCondition
		Orders  []*Order
	}

	// FilterValue is any type value to pass int query filter condition
	FilterValue interface{}

	// QueryBuilder helps with query creation
	QueryBuilder interface {
		Filter(property string, value FilterValue) QueryBuilder
		Order(property string, SortDirection SortDirection)
	}

	// FilterCondition represents a filter operation on a single field
	FilterCondition struct {
		Property  string
		Condition CompareCondition
		Value     FilterValue
	}

	// Order represents a sort operation on a single field
	Order struct {
		Property      string
		SortDirection SortDirection
	}
)

// NewQuery creates a new database query spec. The name is what
// the storage system should use to identify the types, usually
// a table or collection name.
func NewQuery(name string) *Query {
	return &Query{
		Name: name,
	}
}

// Filter adds a filter to the query
func (q *Query) Filter(property string, condition CompareCondition, value FilterValue) *Query {
	filter := NewFilter(property, condition, value)
	q.Filters = append(q.Filters, filter)
	return q
}

// Order adds a sort order to the query
func (q *Query) Order(property string, SortDirection SortDirection) *Query {
	order := NewOrder(property, SortDirection)
	q.Orders = append(q.Orders, order)
	return q
}

// Slice adds a slice operation to the query
func (q *Query) Slice(offset, limit int) *Query {
	q.Offset = offset
	q.Limit = limit
	return q
}

// NewFilter creates a new property filter
func NewFilter(property string, condition CompareCondition, value FilterValue) *FilterCondition {
	return &FilterCondition{
		Property:  property,
		Condition: condition,
		Value:     value,
	}
}

// NewOrder creates a new query order
func NewOrder(property string, SortDirection SortDirection) *Order {
	return &Order{
		Property:      property,
		SortDirection: SortDirection,
	}
}
