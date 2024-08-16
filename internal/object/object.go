package object

import "fmt"

type Type string

const (
	IntegerObj = "INTEGER"
	BooleanObj = "BOOLEAN"
)

type Object interface {
	Type() Type
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i Integer) Type() Type      { return IntegerObj }

type Boolean struct {
	Value bool
}

func (b Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }
