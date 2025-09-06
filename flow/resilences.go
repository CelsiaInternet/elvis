package flow

import (
	"fmt"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/resilience"
)

/**
* AddResilience
* @param flow *Flow
**/
func (s *WorkFlows) addResilience(flow *Flow, ctx et.Json) {
	if flow.Retries > 0 {
		description := fmt.Sprintf("flow: %s,  %s", flow.Name, flow.Description)
		attempt := resilience.AddCustom(flow.Id, flow.Tag, description, flow.Retries, flow.RetryDelay, flow.run, ctx)
		s.Resilience[flow.Id] = attempt
	}
}

/**
* DoneResilience
* @param flow *Flow
**/
func (s *WorkFlows) doneResilience(flow *Flow) {
	if s.Resilience[flow.Id] == nil {
		return
	}

	delete(s.Resilience, flow.Id)
}

/**
* GetResilience
* @param id string
* @return *resilience.Attempt
**/
func (s *WorkFlows) GetResilience(id string) (*resilience.Attempt, error) {
	result := s.Resilience[id]
	if result == nil {
		return nil, fmt.Errorf("resilience not found")
	}

	return result, nil
}
