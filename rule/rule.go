package rule

import (
	"time"

	"github.com/appfirewall/appfirewall/protocol"
)

type Action string

const (
	Allow = Action("allow")
	Deny  = Action("deny")
)

type Duration string

const (
	Once    = Duration("once")
	Restart = Duration("until restart")
	Always  = Duration("always")
)

type Rule struct {
	Created  time.Time
	Updated  time.Time
	Name     string
	Enabled  bool
	Action   Action
	Duration Duration
	Operator Operator
}

func Create(name string, action Action, duration Duration, op Operator) *Rule {
	return &Rule{
		Created:  time.Now(),
		Enabled:  true,
		Name:     name,
		Action:   action,
		Duration: duration,
		Operator: op,
	}
}

func FromAFRule(reply *protocol.AFRule) *Rule {
	operator := NewOperator(
		Type(reply.Operator.Type),
		Operand(reply.Operator.Operand),
		reply.Operator.Data,
		make([]Operator, 0),
	)

	return Create(
		reply.Name,
		Action(reply.Action),
		Duration(reply.Duration),
		operator,
	)
}
