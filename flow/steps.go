package flow

import (
	"fmt"

	"github.com/Knetic/govaluate"
	"github.com/celsiainternet/elvis/et"
)

type Step struct {
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	Stop        bool                           `json:"stop"`
	Expression  string                         `json:"expression"`
	YesGoTo     int                            `json:"yes_go_to"`
	NoGoTo      int                            `json:"no_go_to"`
	flow        *Flow                          `json:"-"`
	fn          FnContext                      `json:"-"`
	rollbacks   FnContext                      `json:"-"`
	expression  *govaluate.EvaluableExpression `json:"-"`
}

/**
* newStep
* @param name, description, expression string, nextIndex int, fn FnContext, stop bool
* @return *Step
**/
func newStep(flow *Flow, name, description string, fn FnContext, stop bool) (*Step, error) {
	result := &Step{
		flow:        flow,
		fn:          fn,
		Name:        name,
		Description: description,
		Stop:        stop,
	}

	return result, nil
}

/**
* run
* @params ctx et.Json
* @return et.Item, error
**/
func (s *Step) run(ctx et.Json) (et.Item, error) {
	s.flow.statusInstance(FlowStatusRunning)
	result, err := s.fn(ctx)
	if err != nil {
		s.flow.statusInstance(FlowStatusFailed)
		return et.Item{}, err
	}

	return result, nil
}

/**
* ToJson
* @return et.Json
**/
func (s *Step) ToJson() et.Json {
	return et.Json{
		"name":        s.Name,
		"description": s.Description,
		"stop":        s.Stop,
		"expression":  s.Expression,
		"yes_go_to":   s.YesGoTo,
		"no_go_to":    s.NoGoTo,
	}
}

/**
* If
* @param expression string, yesGoTo int, noGoTo int
* @return *Step, error
**/
func (s *Step) IfElse(expression string, yesGoTo int, noGoTo int) (*Step, error) {
	s.YesGoTo = yesGoTo
	s.NoGoTo = noGoTo
	if expression != "" {
		evalueExpression, err := govaluate.NewEvaluableExpression(expression)
		if err != nil {
			return s, err
		}

		s.Expression = expression
		s.expression = evalueExpression
	}

	return s, nil
}

/**
* Evaluate
* @param ctx et.Json
* @return bool, error
**/
func (s *Step) Evaluate(ctx et.Json) (bool, error) {
	ok, err := s.expression.Evaluate(ctx)
	if err != nil {
		return false, err
	}

	switch v := ok.(type) {
	case bool:
		return v, nil
	default:
		return false, fmt.Errorf("expression result is not a boolean")
	}
}
