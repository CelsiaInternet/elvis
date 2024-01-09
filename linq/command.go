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
func (c *Linq) Command() (e.Item, error) {
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

	return e.Item{}, nil
}

/**
* Exec
**/
func (c *Linq) commandInsert() (e.Item, error) {
	_, err := c.PrepareInsert()
	if err != nil {
		return e.Item{}, err
	}

	return c.insert()
}

func (c *Linq) commandUpdate() (e.Item, error) {
	current, err := c.PrepareUpdate()
	if err != nil {
		return e.Item{}, err
	}

	return c.update(current)
}

func (c *Linq) commandDelete() (e.Item, error) {
	current, err := c.PrepareDelete()
	if err != nil {
		return e.Item{}, err
	}

	return c.delete(current)
}

func (c *Linq) commandUpsert() (e.Item, error) {
	current, err := c.PrepareUpsert()
	if err != nil {
		return e.Item{}, err
	}

	if current.Ok {
		return c.update(current.Result)
	}

	return c.insert()
}

/**
*
**/
func (c *Linq) Current() (e.Item, error) {
	c.sql = c.SqlCurrent()

	return c.QueryOne()
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
