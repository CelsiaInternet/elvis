package linq

import (
	e "github.com/cgalvisleon/elvis/json"
)

func (c *Linq) Debug() *Linq {
	c.debug = true

	return c
}

/**
* Executors
**/
func (c *Linq) Command() (e.Items, error) {
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

	return e.Items{}, nil
}

func (c *Linq) CommandOne() (e.Item, error) {
	result, err := c.Command()
	if err != nil {
		return e.Item{}, err
	}

	if result.Count == 0 {
		return e.Item{}, nil
	}

	return e.Item{
		Ok:     true,
		Result: result.Result[0],
	}, nil
}

/**
* Exec
**/
func (c *Linq) commandInsert() (e.Items, error) {
	err := c.PrepareInsert()
	if err != nil {
		return e.Items{}, err
	}

	result, err := c.insert()
	if err != nil {
		return e.Items{}, err
	}

	if !result.Ok {
		return e.Items{
			Ok:     false,
			Count:  0,
			Result: []e.Json{},
		}, nil
	}

	return e.Items{
		Ok:     result.Ok,
		Count:  1,
		Result: []e.Json{result.Result},
	}, nil
}

func (c *Linq) commandUpdate() (e.Items, error) {
	var result e.Items = e.Items{}
	currents, err := c.PrepareUpdate()
	if err != nil {
		return e.Items{}, err
	}

	model := c.from[0].model
	for _, current := range currents.Result {
		model.Changue(current, c)
		if c.change {
			item, err := c.update(current)
			if err != nil {
				return e.Items{}, err
			} else {
				result.Result = append(result.Result, item.Result)
				result.Ok = true
				result.Count++
			}
		}
	}

	return result, nil
}

func (c *Linq) commandDelete() (e.Items, error) {
	var result e.Items = e.Items{}
	currents, err := c.PrepareDelete()
	if err != nil {
		return e.Items{}, err
	}

	for _, current := range currents.Result {
		item, err := c.delete(current)
		if err != nil {
			return e.Items{}, err
		} else {
			result.Result = append(result.Result, item.Result)
			result.Ok = true
			result.Count++
		}
	}

	return result, nil
}

func (c *Linq) commandUpsert() (e.Items, error) {
	var result e.Items = e.Items{}
	currents, err := c.PrepareUpsert()
	if err != nil {
		return e.Items{}, err
	}

	if currents.Count == 0 {
		item, err := c.insert()
		if err != nil {
			return e.Items{}, err
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
					return e.Items{}, err
				} else {
					result.Result = append(result.Result, item.Result)
					result.Ok = true
					result.Count++
				}
			}
		}
	}

	return result, nil
}

/**
*
**/
func (c *Linq) Current() (e.Items, error) {
	c.sql = c.SqlCurrent()

	return c.Query()
}

/**
* Basic operation
**/
func (c *Linq) insert() (e.Item, error) {
	model := c.from[0].model

	for _, trigger := range model.BeforeInsert {
		err := trigger(model, nil, c.new, c.data)
		if err != nil {
			return e.Item{}, err
		}
	}

	c.SqlInsert()

	item, err := c.QueryOne()
	if err != nil {
		return e.Item{}, err
	}

	if !item.Ok {
		return item, nil
	}

	c.new = &item.Result

	for _, trigger := range model.AfterInsert {
		err := trigger(model, nil, c.new, c.data)
		if err != nil {
			return e.Item{}, err
		}
	}

	c.Details(&item.Result)

	return item, nil
}

func (c *Linq) update(current e.Json) (e.Item, error) {
	model := c.from[0].model

	for _, trigger := range model.BeforeUpdate {
		err := trigger(model, &current, c.new, c.data)
		if err != nil {
			return e.Item{}, err
		}
	}

	c.SqlUpdate()

	item, err := c.QueryOne()
	if err != nil {
		return e.Item{}, err
	}

	if !item.Ok {
		return item, nil
	}

	c.new = &item.Result

	for _, trigger := range model.AfterUpdate {
		err := trigger(model, &current, c.new, c.data)
		if err != nil {
			return e.Item{}, err
		}
	}

	c.Details(&item.Result)

	return item, nil
}

func (c *Linq) delete(current e.Json) (e.Item, error) {
	model := c.from[0].model

	for _, trigger := range model.BeforeDelete {
		err := trigger(model, &current, nil, c.data)
		if err != nil {
			return e.Item{}, err
		}
	}

	c.SqlDelete()

	item, err := c.QueryOne()
	if err != nil {
		return e.Item{}, err
	}

	if !item.Ok {
		return item, nil
	}

	for _, trigger := range model.AfterDelete {
		err := trigger(model, &current, nil, c.data)
		if err != nil {
			return e.Item{}, err
		}
	}

	return e.Item{
		Ok:     true,
		Result: current,
	}, nil
}
