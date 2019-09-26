package rule

import (
	"fmt"
)

type Type string

const (
	Simple = Type("simple") // Supported for compat (remove?)
	Regexp = Type("regexp") // Supported for compat (remove?)
	List   = Type("list")   // List is the new default (or not even a possibility)
)

const (
	OpList = Operand("list")
)

type Operand string

type opCallback func(value string) bool

type Operator struct {
	Type    Type
	Operand Operand
	Data    string
	List    []Operator
}

func NewOperator(t Type, o Operand, data string, list []Operator) Operator {
	op := Operator{
		Type:    t,
		Operand: o,
		Data:    data,
		List:    list,
	}
	return op
}

func (o *Operator) String() string {
	how := "is"
	if o.Type == Regexp {
		how = "matches"
	}
	return fmt.Sprintf("%s %s '%s'", string(o.Operand), how, string(o.Data))
}
