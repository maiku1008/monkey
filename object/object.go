package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
)

// Object represents any object in the monkey language
type Object interface {
	// Return the ObjectType
	Type() ObjectType
	// Return the value in string form
	Inspect() string
}

// Integer is an object wrapping an integer value according to the monkey language
type Integer struct {
	Value int64
}

var _ Object = (*Integer)(nil)

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// Boolean is an Object wrapping a boolean value according to the monkey language
type Boolean struct {
	Value bool
}

var _ Object = (*Boolean)(nil)

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// Null is an Object wrapping a null value according to the monkey language
type Null struct{}

var _ Object = (*Null)(nil)

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// ReturnValue is an Object wrapping another Object with the return value
type ReturnValue struct {
	Value Object
}

var _ Object = (*ReturnValue)(nil)

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
