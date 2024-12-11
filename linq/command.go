package linq

import (
	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/event"
)

func (c *Linq) Debug() *Linq {
	c.debug = true

	return c
}

/**
* Command
* @return et.Items
* @return error
**/
func (c *Linq) Command() (et.Items, error) {
	if c.Act == ActInsert {
		return c.commandInsert()
	}

	if c.Act == ActUpdate {
		return c.commandUpdate()
	}

	if c.Act == ActUpsert {
		return c.commandUpsert()
	}

	if c.Act == ActDelete {
		return c.commandDelete()
	}

	return et.Items{}, nil
}

/**
* CommandOne
* @return et.Item
* @return error
**/
func (c *Linq) CommandOne() (et.Item, error) {
	result, err := c.Command()
	if err != nil {
		return et.Item{}, err
	}

	if result.Count == 0 {
		return et.Item{}, nil
	}

	return et.Item{
		Ok:     true,
		Result: result.Result[0],
	}, nil
}

/**
* Go
* @return et.Item
* @return error
**/
func (c *Linq) Go() (et.Item, error) {
	return c.CommandOne()
}

/**
* Exec
**/
func (c *Linq) commandInsert() (et.Items, error) {
	currents, err := c.PrepareInsert()
	if err != nil {
		return et.Items{}, err
	}

	if currents.Count > 0 {
		return et.Items{
			Ok:     false,
			Count:  currents.Count,
			Result: currents.Result,
		}, nil
	}

	result, err := c.insert()
	if err != nil {
		return et.Items{}, err
	}

	if !result.Ok {
		return et.Items{}, nil
	}

	return et.Items{
		Ok:     result.Ok,
		Count:  1,
		Result: []et.Json{result.Result},
	}, nil
}

func (c *Linq) commandUpdate() (et.Items, error) {
	var result et.Items = et.Items{}
	currents, err := c.PrepareUpdate()
	if err != nil {
		return et.Items{}, err
	}

	model := c.from[0].model
	for _, current := range currents.Result {
		model.Changue(current, c)
		if c.change {
			item, err := c.update(current)
			if err != nil {
				return et.Items{}, err
			} else {
				result.Result = append(result.Result, item.Result)
				result.Ok = true
				result.Count++
			}
		} else {
			result.Result = append(result.Result, current)
			result.Ok = true
			result.Count++
		}
	}

	return result, nil
}

func (c *Linq) commandDelete() (et.Items, error) {
	var result et.Items = et.Items{}
	currents, err := c.PrepareDelete()
	if err != nil {
		return et.Items{}, err
	}

	for _, current := range currents.Result {
		item, err := c.delete(current)
		if err != nil {
			return et.Items{}, err
		} else {
			result.Result = append(result.Result, item.Result)
			result.Ok = true
			result.Count++
		}
	}

	return result, nil
}

func (c *Linq) commandUpsert() (et.Items, error) {
	var result et.Items = et.Items{}
	currents, err := c.PrepareUpsert()
	if err != nil {
		return et.Items{}, err
	}

	if currents.Count == 0 {
		item, err := c.insert()
		if err != nil {
			return et.Items{}, err
		}

		if item.Ok {
			result.Result = append(result.Result, item.Result)
			result.Ok = true
			result.Count++
		}
	} else {
		model := c.from[0].model
		for _, current := range currents.Result {
			model.Changue(current, c)
			if c.change {
				item, err := c.update(current)
				if err != nil {
					return et.Items{}, err
				} else {
					result.Result = append(result.Result, item.Result)
					result.Ok = true
					result.Count++
				}
			} else {
				result.Result = append(result.Result, current)
				result.Ok = true
				result.Count++
			}
		}
	}

	return result, nil
}

/**
*
**/
func (c *Linq) Current() (et.Items, error) {
	c.sql = c.SqlCurrent()
	result, err := c.query()
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* Basic operation
**/
func (c *Linq) insert() (et.Item, error) {
	model := c.from[0].model

	for _, trigger := range model.BeforeInsert {
		err := trigger(model, nil, c.new, c.data)
		if err != nil {
			return et.Item{}, err
		}
	}

	c.SqlInsert()

	item, err := c.queryOne()
	if err != nil {
		event.Log("error/sql", et.Json{
			"model":  model.Name,
			"action": "insert",
			"sql":    c.sql,
			"error":  err.Error(),
		})
		return et.Item{}, err
	}

	if !item.Ok {
		return item, nil
	}

	new := &item.Result

	for _, trigger := range model.AfterInsert {
		err := trigger(model, nil, new, c.data)
		if err != nil {
			return et.Item{}, err
		}
	}

	c.Details(new)

	return item, nil
}

func (c *Linq) update(current et.Json) (et.Item, error) {
	model := c.from[0].model
	c.idT = current.ValStr("-1", IdTFiled.Low())

	for _, trigger := range model.BeforeUpdate {
		err := trigger(model, &current, c.new, c.data)
		if err != nil {
			return et.Item{}, err
		}
	}

	c.SqlUpdate()

	item, err := c.queryOne()
	if err != nil {
		event.Log("error/sql", et.Json{
			"model":  model.Name,
			"action": "insert",
			"sql":    c.sql,
			"error":  err.Error(),
		})
		return et.Item{}, err
	}

	if !item.Ok {
		return item, nil
	}

	new := &item.Result

	for _, trigger := range model.AfterUpdate {
		err := trigger(model, &current, new, c.data)
		if err != nil {
			return et.Item{}, err
		}
	}

	c.Details(new)

	return item, nil
}

func (c *Linq) delete(current et.Json) (et.Item, error) {
	model := c.from[0].model
	c.idT = current.ValStr("-1", IdTFiled.Low())

	for _, trigger := range model.BeforeDelete {
		err := trigger(model, &current, nil, c.data)
		if err != nil {
			return et.Item{}, err
		}
	}

	c.SqlDelete()

	item, err := c.queryOne()
	if err != nil {
		event.Log("error/sql", et.Json{
			"model":  model.Name,
			"action": "insert",
			"sql":    c.sql,
			"error":  err.Error(),
		})
		return et.Item{}, err
	}

	if !item.Ok {
		return item, nil
	}

	for _, trigger := range model.AfterDelete {
		err := trigger(model, &current, nil, c.data)
		if err != nil {
			return et.Item{}, err
		}
	}

	return et.Item{
		Ok:     true,
		Result: current,
	}, nil
}
