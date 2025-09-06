package flow

import (
	"fmt"

	"github.com/Knetic/govaluate"
	"github.com/celsiainternet/elvis/et"
)

type Step struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Stop        bool      `json:"stop"`
	Expression  string    `json:"expression"`
	YesGoTo     int       `json:"yes_go_to"`
	NoGoTo      int       `json:"no_go_to"`
	fn          FnContext `json:"-"`
	rollbacks   FnContext `json:"-"`
}

/**
* newStep
* @param name, description, expression string, nextIndex int, fn FnContext, stop bool
* @return *Step
**/
func newStep(name, description string, fn FnContext, stop bool) (*Step, error) {
	result := &Step{
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
	result, err := s.fn(ctx)
	if err != nil {
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
func (s *Step) IfElse(expression string, yesGoTo int, noGoTo int) *Step {
	s.YesGoTo = yesGoTo
	s.NoGoTo = noGoTo
	if expression != "" {
		s.Expression = expression
	}

	return s
}

/**
* Evaluate
* @param ctx et.Json
* @return bool, error
**/
func (s *Step) Evaluate(ctx et.Json, instance *Flow) (bool, error) {
	resultError := func(err error) (bool, error) {
		instance.setFailed(fmt.Errorf("error al evaluar expresion:%s, error:%s", s.Expression, err.Error()))
		return false, err
	}

	instance.setStatus(FlowStatusRunning)
	evalueExpression, err := govaluate.NewEvaluableExpression(s.Expression)
	if err != nil {
		return resultError(err)
	}

	ok, err := evalueExpression.Evaluate(ctx)
	if err != nil {
		return resultError(err)
	}

	switch v := ok.(type) {
	case bool:
		return v, nil
	default:
		return resultError(fmt.Errorf("expression result is not a boolean"))
	}
}
